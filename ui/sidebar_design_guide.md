# AI工作流平台 - 侧边栏设计指南

## 设计原则

AI工作流平台的侧边栏设计遵循以下核心原则：

1. **简洁清晰**：简化信息呈现，确保用户能快速找到所需功能
2. **逻辑分组**：相关功能分组组织，提升信息结构的清晰度
3. **视觉层次**：通过设计元素建立清晰的视觉层次，引导用户注意力
4. **一致性**：在整个平台中保持一致的设计语言和交互模式
5. **响应适应**：优雅地适应不同设备尺寸，保持功能的可访问性

## 视觉规范

### 尺寸与空间

- **展开状态宽度**：桌面端 256px (64px [16rem])
- **折叠状态宽度**：桌面端 80px (20px [5rem])
- **移动端宽度**：最小 80px (20px [5rem])
- **内边距**：
  - 水平内边距：16px (4px [1rem])
  - 垂直内边距：16px (4px [1rem])
- **项目高度**：
  - 导航项：48px (12px [3rem])
  - 分类标签：32px (8px [2rem])
  - 用户信息区：64px (16px [4rem])

### 色彩应用

- **背景色**：白色 (#FFFFFF)
- **分割线**：浅灰色 (#E5E7EB)
- **文本颜色**：
  - 主要文本：深灰色 (#1F2937)
  - 次要文本：中灰色 (#6B7280)
  - 高亮文本：靛蓝色 (#4F46E5)
- **图标颜色**：
  - 默认状态：中灰色 (#6B7280)
  - 高亮状态：靛蓝色 (#4F46E5)
- **交互状态**：
  - 悬停背景：浅靛蓝色 (#EFF6FF)
  - 选中背景：浅靛蓝色 (#EFF6FF)
  - 选中边框：无（使用圆角和背景色）

### 排版规范

- **分类标签**：
  - 字体：系统默认无衬线字体
  - 字号：12px (0.75rem)
  - 字重：600 (semibold)
  - 字距：0.05em
  - 大写：全部大写
- **导航项文本**：
  - 字体：系统默认无衬线字体
  - 字号：14px (0.875rem)
  - 字重：500 (medium)
  - 行高：1.5
- **用户名**：
  - 字体：系统默认无衬线字体
  - 字号：14px (0.875rem)
  - 字重：500 (medium)
- **次要文本**：
  - 字体：系统默认无衬线字体
  - 字号：12px (0.75rem)
  - 字重：400 (regular)
  - 颜色：中灰色 (#6B7280)

### 图标与视觉元素

- **导航图标**：
  - 尺寸：20px x 20px
  - 粗细：一致的线条粗细
  - 样式：Font Awesome 图标库
- **徽章**：
  - 尺寸：最小宽度 20px，高度 20px
  - 圆角：9999px (完全圆形)
  - 内边距：水平 8px，垂直 4px
- **用户头像**：
  - 尺寸：32px x 32px
  - 形状：圆形 (border-radius: 100%)
  - 边框：2px 白色边框
- **状态指示点**：
  - 尺寸：8px x 8px
  - 形状：圆形
  - 颜色：
    - 在线：绿色 (#10B981)
    - 离线：灰色 (#9CA3AF)
    - 忙碌：红色 (#EF4444)

## 组件结构

### 侧边栏容器

```html
<aside class="w-20 md:w-64 bg-white shadow-md flex flex-col">
    <!-- 内容 -->
</aside>
```

- 固定高度：100vh 或 100%
- 显示方式：flex 布局，方向为列
- 层叠顺序：适当的 z-index 确保正确覆盖
- 阴影：右侧中等阴影，增强视觉层次

### 侧边栏头部

```html
<div class="p-4 border-b flex items-center justify-center md:justify-start">
    <i class="fas fa-robot text-indigo-600 text-2xl mr-2"></i>
    <span class="hidden md:block text-xl font-bold text-gray-800">AI Agent 平台</span>
</div>
```

- 固定高度：64px
- 边框：下方 1px 分隔线
- Logo：左侧显示应用 Logo
- 标题：在展开状态下显示，折叠时隐藏

### 导航部分

```html
<nav class="flex-1 overflow-y-auto py-4">
    <!-- 导航内容 -->
</nav>
```

- 高度：自适应占满剩余空间
- 滚动：内容过多时可垂直滚动
- 内边距：上下各 16px

### 分类标签

```html
<div class="px-4 py-2">
    <span class="text-xs font-semibold text-gray-500 uppercase tracking-wider">分类名称</span>
</div>
```

- 上下边距：8px
- 左右内边距：16px
- 文本样式：小号、粗体、全大写、增加字距
- 颜色：中灰色，作为视觉分隔而非主要交互元素

### 导航项

```html
<li>
    <a href="page.html" class="group relative flex items-center px-4 py-3 text-gray-700 hover:bg-indigo-50 hover:text-indigo-600 rounded-lg mx-2 transition-all duration-200">
        <i class="fas fa-icon text-lg w-6 text-center"></i>
        <span class="hidden md:block ml-3">菜单名称</span>
        <span class="hidden md:flex ml-auto bg-indigo-100 text-indigo-800 text-xs font-semibold px-2 py-1 rounded-full">徽章</span>
        <!-- 工具提示 -->
        <div class="hidden group-hover:block absolute left-full ml-2 px-2 py-1 bg-gray-800 text-white text-xs rounded whitespace-nowrap z-50">工具提示文本</div>
    </a>
</li>
```

- 容器样式：圆角矩形，悬停/选中时背景色变化
- 布局：Flex 横向排列，对齐居中
- 图标：左侧固定宽度区域
- 文本：中间伸展区域，折叠时隐藏
- 徽章：右侧可选，折叠时隐藏
- 过渡：所有状态变化添加平滑过渡效果
- 工具提示：悬停时在右侧显示

### 选中状态

```html
<a href="current.html" class="group relative flex items-center px-4 py-3 text-indigo-600 sidebar-active hover:bg-indigo-50 hover:text-indigo-600 rounded-lg mx-2 transition-all duration-200">
    <!-- 内容 -->
</a>
```

- 背景色：淡蓝色背景 (#EFF6FF)
- 文本颜色：靛蓝色 (#4F46E5)
- 视觉标识：整个项目设置圆角和背景色，不使用左侧边框

### 用户信息区域

```html
<div class="border-t p-4">
    <a href="profile.html" class="group relative flex items-center text-gray-700 hover:text-indigo-600">
        <div class="relative">
            <img src="user-avatar.jpg" class="w-8 h-8 rounded-full border-2 border-white" alt="用户头像">
            <span class="absolute -top-1 -right-1 w-3.5 h-3.5 bg-green-500 rounded-full border-2 border-white"></span>
        </div>
        <div class="hidden md:block ml-3">
            <p class="text-sm font-medium">用户名</p>
            <p class="text-xs text-gray-500">查看个人资料</p>
        </div>
        <!-- 工具提示 -->
        <div class="hidden group-hover:block absolute left-full ml-2 px-2 py-1 bg-gray-800 text-white text-xs rounded whitespace-nowrap z-50">管理您的个人资料和设置</div>
    </a>
</div>
```

- 分隔：上方 1px 分隔线
- 内边距：16px
- 头像：圆形，带白色边框
- 状态指示：绿色小圆点表示在线状态
- 用户信息：名称和简短说明，折叠时隐藏

### 主题切换区域

```html
<div class="border-t p-4 flex items-center justify-center md:justify-between text-gray-700">
    <span class="hidden md:block text-xs text-gray-500">暗黑模式</span>
    <label class="relative inline-flex items-center cursor-pointer">
        <input type="checkbox" value="" class="sr-only peer">
        <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-2 peer-focus:ring-indigo-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-indigo-600"></div>
    </label>
</div>
```

- 分隔：上方 1px 分隔线
- 内边距：16px
- 布局：两端对齐，移动端居中
- 开关：动画切换效果，遵循系统切换组件规范

## 交互规范

### 鼠标/触摸交互

- **悬停效果**：
  - 导航项背景色淡蓝色渐变变化
  - 文字颜色变为靛蓝色
  - 显示工具提示
  - 过渡时间：200ms
- **点击/触摸效果**：
  - 轻微缩放效果或不透明度变化
  - 背景色短暂加深
  - 触摸反馈适应移动设备
- **选中状态**：
  - 持久的视觉指示（背景色、文字颜色）
  - 区别于临时的悬停效果

### 键盘交互

- **Tab焦点**：
  - 清晰可见的焦点轮廓
  - 遵循自然的导航顺序
  - 焦点状态视觉效果与悬停相似但区分
- **键盘快捷键**：
  - 为主要导航项分配快捷键
  - 遵循平台约定的快捷键规范
  - 提供快捷键提示

### 展开/折叠行为

- **状态切换**：
  - 通过专用控制按钮切换
  - 可选的自动折叠/展开（响应式）
  - 记住用户最后选择的状态
- **过渡动画**：
  - 平滑的宽度变化动画
  - 文字淡入淡出效果
  - 动画时长：300ms

### 响应式行为

- **桌面端**：完整展示所有内容元素
- **平板端**：根据空间考虑折叠
- **移动端**：
  - 默认折叠为图标栏
  - 可临时展开（弹出或推动内容）
  - 可能在小屏幕上转为底部导航

## 辅助功能

### 键盘可访问性

- 所有交互元素可通过键盘访问
- 清晰的焦点状态视觉指示
- 符合逻辑的Tab焦点顺序
- 支持箭头键导航

### 屏幕阅读器支持

- 所有图标设置适当的 aria-label
- 使用语义化HTML结构
- 当前页面使用 aria-current="page"
- 导航分组使用适当的ARIA角色

### 色彩对比度

- 所有文本满足WCAG AA级对比度要求
- 图标和交互元素具有足够的对比度
- 状态指示使用多重视觉提示，不仅依赖颜色
- 分类标签和次要文本满足最低对比度要求

## 实现指南

### HTML结构

遵循以下嵌套层次结构：

```
aside (侧边栏容器)
├── div (侧边栏头部)
├── nav (导航容器)
│   ├── div (分类标签)
│   ├── ul (导航组)
│   │   ├── li (导航项容器)
│   │   │   └── a (导航链接)
│   │   │       ├── i (图标)
│   │   │       ├── span (文本)
│   │   │       ├── span (徽章)
│   │   │       └── div (工具提示)
│   │   └── ... (更多导航项)
│   └── ... (更多分类和导航组)
├── div (用户信息区)
└── div (其他底部元素，如主题切换)
```

### CSS类命名

建议使用Tailwind CSS类名，保持一致性：

- 容器：`sidebar`, `sidebar-header`, `sidebar-nav`, etc.
- 状态：`sidebar-active`, `hover:bg-indigo-50`, etc.
- 工具提示：`group`, `group-hover:block`, etc.
- 响应式：`hidden`, `md:block`, `lg:flex`, etc.

### JavaScript交互

为侧边栏实现以下交互功能：

- **折叠/展开切换**：
  ```javascript
  // 切换侧边栏状态
  document.querySelector('.sidebar-toggle').addEventListener('click', function() {
    document.querySelector('.sidebar').classList.toggle('sidebar-collapsed');
  });
  ```

- **响应式调整**：
  ```javascript
  // 根据窗口大小自动调整
  function adjustSidebar() {
    const sidebar = document.querySelector('.sidebar');
    if (window.innerWidth < 768) {
      sidebar.classList.add('sidebar-collapsed');
    } else {
      sidebar.classList.remove('sidebar-collapsed');
    }
  }
  
  window.addEventListener('resize', adjustSidebar);
  adjustSidebar(); // 初始调整
  ```

- **当前页面高亮**：
  ```javascript
  // 高亮当前页面对应的导航项
  document.addEventListener('DOMContentLoaded', function() {
    const currentPath = window.location.pathname;
    const navLinks = document.querySelectorAll('.sidebar nav a');
    
    navLinks.forEach(link => {
      if (link.getAttribute('href') === currentPath.split('/').pop()) {
        link.classList.add('sidebar-active');
      } else {
        link.classList.remove('sidebar-active');
      }
    });
  });
  ```

### 性能注意事项

- 避免过多的JavaScript计算，优先使用CSS实现交互效果
- 使用CSS过渡而非JavaScript动画
- 图标使用CSS图标字体或内联SVG，避免额外的网络请求
- 考虑在移动设备上延迟加载非关键导航项

## 设计变体

### 亮色主题

- 背景色：白色 (#FFFFFF)
- 文本：深灰色 (#1F2937)
- 高亮：靛蓝色 (#4F46E5)
- 分隔线：浅灰色 (#E5E7EB)

### 暗色主题

- 背景色：深灰色 (#1F2937)
- 文本：浅灰色 (#F3F4F6)
- 高亮：靛蓝色 (#818CF8)
- 分隔线：中灰色 (#4B5563)

### 迷你变体

- 永久折叠状态
- 只显示图标
- 悬停时显示工具提示
- 适用于需要更多工作区空间的场景

### 顶部导航变体

- 水平布局，位于页面顶部
- 导航项左右排列
- 下拉菜单取代分组
- 适用于移动设备或特定布局需求

## 最佳实践

1. **简化导航层级**：避免过深的导航层级，保持扁平化结构
2. **分组优化**：相关功能分组，但每组不超过7项
3. **视觉区分**：确保选中状态有明确的视觉差异
4. **一致交互**：保持与平台其他部分一致的交互模式
5. **易于扩展**：设计应考虑未来添加新导航项的可能
6. **性能优先**：避免过于复杂的动画和效果影响性能
7. **响应式优先**：从移动设备开始设计，逐步扩展到大屏幕

## 错误与边界情况

1. **过长的菜单文本**：设置最大宽度和溢出处理
2. **过多的导航项**：确保滚动行为正常，考虑分组或折叠
3. **徽章数值过大**：设置最大值显示（如"99+"）
4. **无JavaScript降级**：确保基本功能在禁用JavaScript时仍可用
5. **图片加载失败**：为头像和图标提供后备方案
6. **屏幕尺寸极端情况**：测试极小和极大屏幕下的表现

## 测试清单

- [ ] 所有链接正确指向目标页面
- [ ] 响应式布局在所有目标设备尺寸下正常工作
- [ ] 键盘导航功能完整，Tab顺序合理
- [ ] 屏幕阅读器可以正确解读所有导航元素
- [ ] 色彩对比度满足可访问性要求
- [ ] 状态变化有合适的视觉反馈
- [ ] 动画和过渡效果流畅自然
- [ ] 暗色模式下所有元素正常显示
- [ ] 所有图标加载正确且含义明确
- [ ] 折叠/展开功能正常工作 