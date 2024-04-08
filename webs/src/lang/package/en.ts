export default {
  // 路由国际化
  route: {
    dashboard: "Dashboard",
    document: "Document",
    userset: "ChangePassword",
    system:"system management",
    nodelist:"Node Management",
    sublist:"Sub Management",
    subcription:"Subcription",
  },
  // 登录页面国际化
  login: {
    username: "Username",
    password: "Password",
    login: "Login",
    captchaCode: "Verify Code",
    capsLock: "Caps Lock is On",
    message: {
      username: {
        required: "Please enter Username",
      },
      password: {
        required: "Please enter Password",
        min: "The password can not be less than 6 digits",
      },
      captchaCode: {
        required: "Please enter Verify Code",
      },
    },
  },
  // 重置密码页面国际化
  userset:{
    title: "Change Password",
    newUsername: "New Username",
    newPassword: "New Password",
    message: {
      title:"Prompt",
      xx1:"Username or password cannot be empty",
      xx2: "The password length cannot be less than 6 digits",
      xx3:"Are you sure you want to reset the password",
      xx4:"Password reset successful, new password is:",
    },
  },
  // 导航栏国际化
  navbar: {
    dashboard: "Dashboard",
    logout: "Logout",
    userset: "ChangePassword",
  },
  sizeSelect: {
    tooltip: "Layout Size",
    default: "Default",
    large: "Large",
    small: "Small",
    message: {
      success: "Switch Layout Size Successful!",
    },
  },
  langSelect: {
    message: {
      success: "Switch Language Successful!",
    },
  },
  settings: {
    project: "Project Settings",
    theme: "Theme",
    interface: "Interface",
    navigation: "Navigation",
    themeColor: "Theme Color",
    tagsView: "Tags View",
    fixedHeader: "Fixed Header",
    sidebarLogo: "Sidebar Logo",
    watermark: "Watermark",
  },
};
