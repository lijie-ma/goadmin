<template>
  <div class="login-container">
    <div class="login-box">
      <h2>管理系统登录</h2>
      <form @submit.prevent="handleLogin">
        <div class="form-item">
          <input
            type="text"
            v-model="formData.username"
            placeholder="请输入用户名"
            required
          >
        </div>
        <div class="form-item">
          <input
            type="password"
            v-model="formData.password"
            placeholder="请输入密码"
            required
          >
        </div>
        <div class="form-item remember-me">
          <label>
            <input
              type="checkbox"
              v-model="formData.remember"
            > 记住我
          </label>
        </div>
        <div class="form-item">
          <button type="submit" :disabled="loading">
            {{ loading ? '登录中...' : '登录' }}
          </button>
        </div>
      </form>
    </div>

    <!-- 滑动验证码弹框 -->
    <div class="captcha-modal" v-if="showCaptcha">
      <div class="captcha-dialog">
        <div class="captcha-header">
          <span class="captcha-title">请完成安全验证</span>
          <button class="close-btn" @click="closeCaptcha">×</button>
        </div>

        <div class="canvas-wrap" id="canvasWrap">
          <canvas ref="bgCanvas" id="bg" width="300" height="220"></canvas>
          <canvas ref="pieceCanvas" id="piece" width="300" height="220"></canvas>
          <div class="mask" v-if="captchaLoading">加载中...</div>
        </div>

        <div class="slider-area">
          <div class="track" ref="track">
            <div class="progress" ref="progress"></div>
            <div
              class="thumb"
              ref="thumb"
              @mousedown="pointerDown"
              @touchstart="pointerDown"
            >
              <div class="bar"></div>
            </div>
          </div>
        </div>

        <div class="controls">
          <button class="refresh-btn" @click="loadCaptcha">刷新</button>
        </div>
        <div class="msg" :class="msgClass">{{ captchaMsg }}</div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onUnmounted } from 'vue'

// 表单数据
const formData = reactive({
  username: '',
  password: '',
  token: '',
  remember: false
})

// 加载状态
const loading = ref(false)

// 验证码相关状态
const showCaptcha = ref(false)
const captchaLoading = ref(false)
const captchaMsg = ref('')
const msgClass = ref('')

// DOM引用
const bgCanvas = ref(null)
const pieceCanvas = ref(null)
const thumb = ref(null)
const track = ref(null)
const progress = ref(null)

// 验证码状态
const captchaState = reactive({
  token: null,
  pieceX: 0,
  pieceY: 0,
  pieceW: 60,
  dragging: false,
  startX: 0,
  thumbStartLeft: 0,
  maxMove: 0,
  pieceData: null
})

// 验证码API配置
const captchaConfig = {
  newUrl: '/api/admin/v1/captcha/generate',
  verifyUrl: '/api/admin/v1/captcha/check',
  checkUrl: '/api/admin/v1/setting/get_captcha_switch'  // 检查是否需要验证码的接口
}

// 关闭验证码弹框
const closeCaptcha = () => {
  showCaptcha.value = false
  resetSlider()
}

// 重置滑块
const resetSlider = () => {
  if (progress.value && thumb.value) {
    progress.value.style.width = '0'
    thumb.value.style.transform = 'translateX(0px)'
    if (pieceCanvas.value) {
      pieceCanvas.value.style.transform = 'translateX(0px)'
    }
  }
}

// 设置验证码消息
const setMsg = (text, type = '') => {
  captchaMsg.value = text
  msgClass.value = type ? `msg ${type}` : 'msg'
}

// 计算滑块偏移对应的拼图偏移
const thumbToOffset = (thumbX) => {
  if (!bgCanvas.value || !captchaState.pieceW) return 0
  const maxOffset = bgCanvas.value.width - captchaState.pieceW
  return Math.round((thumbX / captchaState.maxMove) * maxOffset)
}

// 更新拼图位置
const updatePiece = (thumbX) => {
  if (!pieceCanvas.value || !captchaState.pieceData) return
  const offset = thumbToOffset(thumbX)
  pieceCanvas.value.style.transform = `translateX(${offset}px)`
}

// 加载验证码
const loadCaptcha = async () => {
  setMsg('')
  captchaLoading.value = true
  resetSlider()

  try {
    const response = await fetch(captchaConfig.newUrl, {
      cache: 'no-store'
    })
    if (!response.ok) throw new Error('网络错误')

    const res = await response.json()
    captchaState.token = res.data.key
    const bgSrc = res.data.image_base64
    const pieceSrc = res.data.tile_base64
    captchaState.pieceX = res.data.tile_x
    captchaState.pieceY = res.data.tile_y
    captchaState.pieceW = res.data.tile_width || 60

    await new Promise((resolve, reject) => {
      const bgImg = new Image()
      bgImg.onload = () => {
        if (!bgCanvas.value) return
        const ctx = bgCanvas.value.getContext('2d')
        ctx.drawImage(bgImg, 0, 0, bgCanvas.value.width, bgCanvas.value.height)
        ctx.clearRect(
          captchaState.pieceX,
          captchaState.pieceY,
          captchaState.pieceW,
          captchaState.pieceW
        )

        if (pieceSrc) {
          const pieceImg = new Image()
          pieceImg.onload = () => {
            if (!pieceCanvas.value) return
            const pctx = pieceCanvas.value.getContext('2d')
            pctx.drawImage(
              pieceImg,
              0,
              captchaState.pieceY,
              captchaState.pieceW,
              captchaState.pieceW
            )

            // 保存拼图数据用于拖动
            const temp = document.createElement('canvas')
            temp.width = captchaState.pieceW
            temp.height = captchaState.pieceW
            const tctx = temp.getContext('2d')
            tctx.drawImage(pieceImg, 0, 0, captchaState.pieceW, captchaState.pieceW)
            captchaState.pieceData = tctx.getImageData(0, 0, captchaState.pieceW, captchaState.pieceW)

            if (track.value && thumb.value) {
              const trackRect = track.value.getBoundingClientRect()
              const thumbRect = thumb.value.getBoundingClientRect()
              captchaState.maxMove = trackRect.width - thumbRect.width
            }
            resolve()
          }
          pieceImg.onerror = reject
          pieceImg.src = pieceSrc
        } else {
          resolve()
        }
      }
      bgImg.onerror = reject
      bgImg.src = bgSrc
    })

    captchaLoading.value = false
  } catch (err) {
    console.error('加载验证码失败:', err)
    setMsg('加载验证码失败', 'fail')
    captchaLoading.value = false
  }
}

// 处理鼠标/触摸按下
const pointerDown = (e) => {
  e.preventDefault()
  captchaState.dragging = true
  captchaState.startX = e.touches ? e.touches[0].clientX : e.clientX
  const transform = thumb.value?.style.transform || 'translateX(0px)'
  const m = transform.match(/translateX\(([-\d.]+)px\)/)
  captchaState.thumbStartLeft = m ? parseFloat(m[1]) : 0
  setMsg('')

  // 添加移动和释放事件监听
  document.addEventListener('mousemove', pointerMove)
  document.addEventListener('touchmove', pointerMove, { passive: false })
  document.addEventListener('mouseup', pointerUp)
  document.addEventListener('touchend', pointerUp)
}

// 处理鼠标/触摸移动
const pointerMove = (e) => {
  if (!captchaState.dragging) return

  const clientX = e.touches ? e.touches[0].clientX : e.clientX
  const dx = clientX - captchaState.startX
  let newLeft = captchaState.thumbStartLeft + dx

  // 限制滑动范围
  if (newLeft < 0) newLeft = 0
  if (newLeft > captchaState.maxMove) newLeft = captchaState.maxMove

  // 更新滑块位置
  if (thumb.value) {
    thumb.value.style.transform = `translateX(${newLeft}px)`
  }

  // 更新进度条
  if (progress.value && thumb.value) {
    progress.value.style.width = `${newLeft + thumb.value.offsetWidth/2}px`
  }

  // 更新拼图位置
  updatePiece(newLeft)
}

// 处理鼠标/触摸释放
const pointerUp = async (e) => {
  if (!captchaState.dragging) return
  captchaState.dragging = false

  // 移除事件监听
  document.removeEventListener('mousemove', pointerMove)
  document.removeEventListener('touchmove', pointerMove)
  document.removeEventListener('mouseup', pointerUp)
  document.removeEventListener('touchend', pointerUp)

  // 获取最终位置
  const transform = thumb.value?.style.transform || 'translateX(0px)'
  const m = transform.match(/translateX\(([-\d.]+)px\)/)
  const left = m ? parseFloat(m[1]) : 0
  const offset = thumbToOffset(left)

  setMsg('验证中...')
  try {
    const response = await fetch(captchaConfig.verifyUrl, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        key: captchaState.token,
        x: offset,
        y: captchaState.pieceY
      })
    })

    const result = await response.json()
    if (result.code === 200) {
      setMsg('验证通过', 'success')
      if (track.value && progress.value) {
        progress.value.style.width = track.value.getBoundingClientRect().width + 'px'
      }
      setTimeout((result) => {
        showCaptcha.value = false
        handleLoginSubmit(result.data)
      }, 500, result)
    } else {
      setMsg(result.message || '验证失败', 'fail')
      resetSlider()
    }
  } catch (err) {
    console.error('验证请求失败:', err)
    setMsg('验证请求失败', 'fail')
    resetSlider()
  }
}

// 处理登录表单提交
const handleLoginSubmit = async (captchaData) => {
  try {
    formData.token = captchaData?.token
    loading.value = true
    const response = await fetch('/api/admin/v1/user/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(formData)
    })

    if (!response.ok) {
      throw new Error('登录失败')
    }

    const data = await response.json()
    console.log('登录成功', data)
    // 这里可以处理登录成功后的操作，如保存token等
  } catch (error) {
    console.error('登录失败:', error)
    alert('登录失败，请重试')
  } finally {
    loading.value = false
  }
}

// 处理登录按钮点击
const handleLogin = async () => {
  try {
    // 先检查是否需要验证码
    const checkResponse = await fetch(captchaConfig.checkUrl, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json'
      }
    })

    if (checkResponse.status !== 200) {
      throw new Error('检查验证码状态失败')
    }

    const checkResult = await checkResponse.json()

    // 判断是否需要验证码
    if (checkResult.data.admin === 1) {
      // 需要验证码，显示验证码弹框
      showCaptcha.value = true
      await loadCaptcha()
    } else {
      // 不需要验证码，直接登录
      await handleLoginSubmit()
    }
  } catch (error) {
    console.error('检查验证码状态失败:', error)
    // 出错时默认显示验证码
    showCaptcha.value = true
    await loadCaptcha()
  }
}

// 组件卸载时清理事件监听
onUnmounted(() => {
  document.removeEventListener('mousemove', pointerMove)
  document.removeEventListener('touchmove', pointerMove)
  document.removeEventListener('mouseup', pointerUp)
  document.removeEventListener('touchend', pointerUp)
})
</script>

<style scoped>
/* 登录表单样式 */
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background-color: #f0f2f5;
}

.login-box {
  width: 100%;
  max-width: 400px;
  padding: 40px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

h2 {
  text-align: center;
  margin-bottom: 30px;
  color: #1a1a1a;
}

.form-item {
  margin-bottom: 20px;
}

input[type="text"],
input[type="password"] {
  width: 100%;
  padding: 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
  transition: border-color 0.3s;
}

input[type="text"]:focus,
input[type="password"]:focus {
  border-color: #1890ff;
  outline: none;
}

.remember-me {
  display: flex;
  align-items: center;
  font-size: 14px;
  color: #666;
}

button {
  width: 100%;
  padding: 12px;
  background-color: #1890ff;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 16px;
  cursor: pointer;
  transition: background-color 0.3s;
}

button:hover {
  background-color: #40a9ff;
}

button:disabled {
  background-color: #ccc;
  cursor: not-allowed;
}

/* 验证码弹框样式 */
.captcha-modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.captcha-dialog {
  width: 340px;
  background: #fff;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.15);
}

.captcha-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.captcha-title {
  font-size: 16px;
  font-weight: 500;
  color: #333;
}

.close-btn {
  background: none;
  border: none;
  font-size: 20px;
  color: #999;
  cursor: pointer;
  padding: 0;
  width: 24px;
  height: 24px;
  line-height: 24px;
}

.close-btn:hover {
  color: #666;
}

.canvas-wrap {
  position: relative;
  width: 300px;
  height: 220px;
  margin: 0 auto;
  background: #e9eef5;
  border-radius: 6px;
  overflow: hidden;
}

canvas {
  display: block;
  border-radius: 6px;
}

#piece {
  position: absolute;
  top: 0;
  left: 0;
  pointer-events: none;
  transition: transform 0.05s linear;
}

.mask {
  position: absolute;
  left: 0;
  top: 0;
  right: 0;
  bottom: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  color: #666;
  background: rgba(255, 255, 255, 0.8);
}

.slider-area {
  width: 300px;
  margin: 12px auto 0;
  height: 40px;
  display: flex;
  align-items: center;
  background: #f3f6f9;
  border-radius: 20px;
  padding: 6px;
  box-sizing: border-box;
}

.track {
  position: relative;
  flex: 1;
  height: 100%;
  background: #fff;
  border-radius: 14px;
  box-shadow: inset 0 0 0 1px rgba(0, 0, 0, 0.04);
  display: flex;
  align-items: center;
}

.progress {
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  background: linear-gradient(90deg, #2d8cf0, #4aa3ff);
  border-radius: 14px;
  width: 0;
  transition: none;
}

.thumb {
  width: 40px;
  height: 28px;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.12);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  user-select: none;
}

.thumb .bar {
  width: 10px;
  height: 10px;
  background: #2d8cf0;
  border-radius: 2px;
}

.controls {
  display: flex;
  gap: 12px;
  margin-top: 10px;
  align-items: center;
  justify-content: center;
}

.refresh-btn {
  width: auto;
  padding: 8px 16px;
}

.msg {
  text-align: center;
  margin-top: 8px;
  color: #666;
  font-size: 13px;
  min-height: 18px;
}

.msg.success {
  color: #52c41a;
}

.msg.fail {
  color: #ff4d4f;
}
</style>
