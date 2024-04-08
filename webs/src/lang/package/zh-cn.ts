export default {
  // 路由国际化
  route: {
    dashboard: "首页",
    document: "项目文档",
    userset: "修改密码",
    system:"系统管理",
    nodelist:"节点管理",
    sublist:"订阅管理",
    subcription:"节点和订阅",
  },
  // 登录页面国际化
  login: {
    username: "用户名",
    password: "密码",
    login: "登 录",
    captchaCode: "验证码",
    capsLock: "大写锁定已打开",
    message: {
      username: {
        required: "请输入用户名",
      },
      password: {
        required: "请输入密码",
        min: "密码不能少于6位",
      },
      captchaCode: {
        required: "请输入验证码",
      },
    },
  },
    // 重置密码页面国际化
    userset:{
      title: "修改密码",
      newUsername: "新账号",
      newPassword: "新密码",
      message: {
        title:"提示",
        xx1:"账号或密码不能为空",
        xx2: "密码长度不能小于6位",
        xx3:"你确定要重置密码吗",
        xx4:"密码重置成功，新密码是：",
      },
    },
  // 导航栏国际化
  navbar: {
    dashboard: "首页",
    logout: "注销登出",
    userset: "修改密码",
  },
  sizeSelect: {
    tooltip: "布局大小",
    default: "默认",
    large: "大型",
    small: "小型",
    message: {
      success: "切换布局大小成功！",
    },
  },
  langSelect: {
    message: {
      success: "切换语言成功！",
    },
  },
  settings: {
    project: "项目配置",
    theme: "主题设置",
    interface: "界面设置",
    navigation: "导航设置",
    themeColor: "主题颜色",
    tagsView: "开启 Tags-View",
    fixedHeader: "固定 Header",
    sidebarLogo: "侧边栏 Logo",
    watermark: "开启水印",
  },
};
