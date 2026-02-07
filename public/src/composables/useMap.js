import { ref, nextTick } from 'vue'
import { ElMessage } from 'element-plus'

export function useMap(serviceSettings, amapKey) {
  // 新增地图相关
  let addMap = null
  let addMarker = null
  let addPlaceSearch = null
  let addGeocoder = null
  const mapLoaded = ref(false)

  // 编辑地图相关
  let editMap = null
  let editMarker = null
  let editPlaceSearch = null
  let editGeocoder = null

  // 加载高德地图API
  const loadAMapScript = () => {
    return new Promise((resolve, reject) => {
      if (window.AMap) {
        resolve()
        return
      }

      // 检查安全配置是否已设置
      if (!window._AMapSecurityConfig) {
        console.error('AMap security config not set')
        reject(new Error('AMap security config not set'))
        return
      }

      const script = document.createElement('script')
      script.type = 'text/javascript'
      script.src = `https://webapi.amap.com/maps?v=2.0&key=${amapKey.value}&plugin=AMap.PlaceSearch,AMap.Geocoder`
      script.onload = () => {
        mapLoaded.value = true
        resolve()
      }
      script.onerror = () => {
        reject(new Error('Failed to load AMap script'))
      }
      document.head.appendChild(script)
    })
  }

  // 初始化新增地图
  const initAddMap = async (initialCityParam = '', onMapClick) => {
    try {
      await loadAMapScript()
      await nextTick()

      const mapContainer = document.getElementById('addMapContainer')
      if (!mapContainer) {
        console.error('Map container not found')
        return
      }

      const initialCity =
        initialCityParam ||
        serviceSettings?.region ||
        '全国'

      console.log('Initializing map with city:', initialCity, 'serviceSettings:', serviceSettings)

      // 创建地图实例
      addMap = new AMap.Map('addMapContainer', {
        zoom: 11,
        viewMode: '2D',
        resizeEnable: true
      })

      // 动态加载控件插件
      AMap.plugin(['AMap.ControlBar', 'AMap.Scale'], () => {
        addMap.addControl(
          new AMap.ControlBar({
            position: { right: '10px', top: '10px' }
          })
        )
        addMap.addControl(new AMap.Scale())
      })

      // 初始化搜索与地理编码服务
      addPlaceSearch = new AMap.PlaceSearch({
        city: initialCity,
        pageSize: 5,
        pageIndex: 1,
        extensions: 'all'
      })

      addGeocoder = new AMap.Geocoder({
        city: initialCity
      })

      // 地图点击事件
      addMap.on('click', (e) => {
        onMapClick(e.lnglat)
      })

      // 设置地图中心
      if (initialCity && initialCity !== '全国') {
        setMapCenterByCity(initialCity)
      } else {
        console.log('No specific city, using default China center')
        addMap.setZoom(4)
        addMap.setCenter([104.195397, 35.86166])
      }
    } catch (error) {
      console.error('Failed to initialize map:', error)
      ElMessage.error('地图加载失败')
    }
  }

  // 根据城市设置地图中心
  const setMapCenterByCity = (city) => {
    if (!addGeocoder || !city) {
      console.log('Cannot set map center: geocoder or city missing')
      return
    }

    addGeocoder.getLocation(city, (status, result) => {
      if (status === 'complete' && result.geocodes?.length > 0) {
        const location = result.geocodes[0].location
        addMap.setCenter([location.lng, location.lat])
        addMap.setZoom(11)
      } else {
        console.warn(`City not found: ${city}`)
        ElMessage.warning(`未找到城市: ${city}`)
        addMap.setZoom(4)
        addMap.setCenter([104.195397, 35.86166])
      }
    })
  }

  // 搜索地点
  const handleSearchLocation = (searchKeyword, onLocationFound) => {
    if (!searchKeyword || !addPlaceSearch) {
      return
    }

    addPlaceSearch.search(searchKeyword, (status, result) => {
      if (status === 'complete' && result.poiList.pois.length > 0) {
        const poi = result.poiList.pois[0]
        const location = poi.location

        // 设置地图中心和标记
        addMap.setCenter([location.lng, location.lat])
        addMap.setZoom(15)

        // 更新标记
        updateMarker(location.lng, location.lat)

        // 回调处理位置找到
        onLocationFound(poi, location)
      } else {
        ElMessage.warning('未找到该位置')
      }
    })
  }

  // 处理地图点击
  const handleMapClick = (lnglat, onLocationUpdate) => {
    updateMarker(lnglat.lng, lnglat.lat)

    // 反向地理编码获取地址
    if (addGeocoder) {
      addGeocoder.getAddress(lnglat, (status, result) => {
        if (status === 'complete' && result.regeocode) {
          onLocationUpdate(result.regeocode.formattedAddress, lnglat.lng, lnglat.lat)
        }
      })
    } else {
      onLocationUpdate('', lnglat.lng, lnglat.lat)
    }
  }

  // 更新地图标记
  const updateMarker = (lng, lat) => {
    if (addMarker) {
      addMarker.setPosition([lng, lat])
    } else {
      addMarker = new AMap.Marker({
        position: [lng, lat],
        map: addMap,
        draggable: true
      })

      // 标记拖拽事件
      addMarker.on('dragend', (e) => {
        const position = e.target.getPosition()
        handleMapClick(position, () => {})
      })
    }
  }

  // 初始化编辑地图
  const initEditMap = async (initialPosition, onMapClick) => {
    try {
      await loadAMapScript()
      await nextTick()

      const mapContainer = document.getElementById('editMapContainer')
      if (!mapContainer) {
        console.error('Edit map container not found')
        return
      }

      const initialCity = serviceSettings?.region || '全国'

      // 创建地图实例
      editMap = new AMap.Map('editMapContainer', {
        zoom: 11,
        viewMode: '2D',
        resizeEnable: true
      })

      // 动态加载控件插件
      AMap.plugin(['AMap.ControlBar', 'AMap.Scale'], () => {
        editMap.addControl(
          new AMap.ControlBar({
            position: { right: '10px', top: '10px' }
          })
        )
        editMap.addControl(new AMap.Scale())
      })

      // 初始化搜索与地理编码服务
      editPlaceSearch = new AMap.PlaceSearch({
        city: initialCity,
        pageSize: 5,
        pageIndex: 1,
        extensions: 'all'
      })

      editGeocoder = new AMap.Geocoder({
        city: initialCity
      })

      // 地图点击事件
      editMap.on('click', (e) => {
        onMapClick(e.lnglat)
      })

      // 设置地图中心到当前位置
      if (initialPosition?.longitude && initialPosition?.latitude) {
        editMap.setCenter([initialPosition.longitude, initialPosition.latitude])
        editMap.setZoom(15)
        updateEditMarker(initialPosition.longitude, initialPosition.latitude)
      } else if (initialCity && initialCity !== '全国') {
        setEditMapCenterByCity(initialCity)
      } else {
        editMap.setZoom(4)
        editMap.setCenter([104.195397, 35.86166])
      }
    } catch (error) {
      console.error('Failed to initialize edit map:', error)
      ElMessage.error('地图加载失败')
    }
  }

  // 根据城市设置编辑地图中心
  const setEditMapCenterByCity = (city) => {
    if (!editGeocoder || !city) {
      console.log('Cannot set edit map center: geocoder or city missing')
      return
    }

    editGeocoder.getLocation(city, (status, result) => {
      if (status === 'complete' && result.geocodes?.length > 0) {
        const location = result.geocodes[0].location
        editMap.setCenter([location.lng, location.lat])
        editMap.setZoom(11)
      } else {
        console.warn(`City not found: ${city}`)
        editMap.setZoom(4)
        editMap.setCenter([104.195397, 35.86166])
      }
    })
  }

  // 编辑地图搜索地点
  const handleEditSearchLocation = (searchKeyword, onLocationFound) => {
    if (!searchKeyword || !editPlaceSearch) {
      return
    }

    editPlaceSearch.search(searchKeyword, (status, result) => {
      if (status === 'complete' && result.poiList.pois.length > 0) {
        const poi = result.poiList.pois[0]
        const location = poi.location

        // 设置地图中心和标记
        editMap.setCenter([location.lng, location.lat])
        editMap.setZoom(15)

        // 更新标记
        updateEditMarker(location.lng, location.lat)

        // 回调处理位置找到
        onLocationFound(poi, location)
      } else {
        ElMessage.warning('未找到该位置')
      }
    })
  }

  // 处理编辑地图点击
  const handleEditMapClick = (lnglat, onLocationUpdate) => {
    updateEditMarker(lnglat.lng, lnglat.lat)

    // 反向地理编码获取地址
    if (editGeocoder) {
      editGeocoder.getAddress(lnglat, (status, result) => {
        if (status === 'complete' && result.regeocode) {
          onLocationUpdate(result.regeocode.formattedAddress, lnglat.lng, lnglat.lat)
        }
      })
    } else {
      onLocationUpdate('', lnglat.lng, lnglat.lat)
    }
  }

  // 更新编辑地图标记
  const updateEditMarker = (lng, lat) => {
    if (editMarker) {
      editMarker.setPosition([lng, lat])
    } else {
      editMarker = new AMap.Marker({
        position: [lng, lat],
        map: editMap,
        draggable: true
      })

      // 标记拖拽事件
      editMarker.on('dragend', (e) => {
        const position = e.target.getPosition()
        handleEditMapClick(position, () => {})
      })
    }
  }

  // 清理新增地图资源
  const destroyAddMap = () => {
    if (addMap) {
      addMap.destroy()
      addMap = null
    }
    addMarker = null
    addPlaceSearch = null
    addGeocoder = null
  }

  // 清理编辑地图资源
  const destroyEditMap = () => {
    if (editMap) {
      editMap.destroy()
      editMap = null
    }
    editMarker = null
    editPlaceSearch = null
    editGeocoder = null
  }

  return {
    mapLoaded,
    initAddMap,
    handleSearchLocation,
    handleMapClick,
    initEditMap,
    handleEditSearchLocation,
    handleEditMapClick,
    destroyAddMap,
    destroyEditMap
  }
}