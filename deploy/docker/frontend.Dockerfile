# 构建阶段
FROM node:20-alpine as builder

# 设置工作目录
WORKDIR /app

# 复制package.json和package-lock.json
COPY public/src/package*.json ./

# 安装依赖
RUN npm install

# 复制源代码
COPY public/src .

# 构建项目
RUN npm run build

# 生产阶段
FROM nginx:alpine

# 复制构建产物到nginx目录
COPY --from=builder /app/dist /usr/share/nginx/html

# 复制nginx配置文件
COPY deploy/docker/nginx.conf /etc/nginx/conf.d/default.conf

# 暴露80端口
EXPOSE 80

# 启动nginx
CMD ["nginx", "-g", "daemon off;"]
