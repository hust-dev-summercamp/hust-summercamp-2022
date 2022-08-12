import { loginUser, getProfile, registerUser } from '@/api/sys'
import md5 from 'md5'
import { message } from '@/libs'
import { LOGIN_TYPE_OAUTH_NO_REGISTER_CODE } from '@/constants'

export default {
  namespaced: true,
  state: () => ({
    // 登录之后的 token
    token: '',
    // 用户名
    username: '',
    // 获取用户信息
    userInfo: {}
  }),
  mutations: {
    /**
     * 保存 token
     */
    setToken(state, newToken) {
      state.token = newToken
    },
    setUsername(state, newName) {
      state.username = newName
    },
    /**
     * 保存用户信息
     */
    setUserInfo(state, newInfo) {
      state.userInfo = newInfo
    }
  },
  actions: {
    /**
     * 注册
     */
    async register(context, payload) {
      // 注册
      return await registerUser({
        ...payload,
      })
    },
    /**
     * 登录
     */
    async login(context, payload) {
      const { password } = payload
      const data = await loginUser({
        ...payload,
      })

      context.commit('setToken', data.token)
      context.commit('setUsername', data.username)
    },
    /**
     * 获取用户信息
     */
    // async profile(context) {
    //   const data = await getProfile()
    //   context.commit('setUserInfo', data)
    //   // 欢迎
    //   message(
    //     'success',
    //     `欢迎您 ${
    //       data.vipLevel
    //         ? '尊贵的 VIP' + data.vipLevel + ' 用户 ' + data.nickname
    //         : data.nickname
    //     } `,
    //     6000
    //   )
    // },
    /**
     * 退出登录
     */
    logout(context) {
      context.commit('setToken', '')
      context.commit('setUserInfo', {})
      // 退出登录之后，重新刷新下页面，因为对于前台项目而言，用户是否登录（是否为 VIP）看到的数据可能不同
      location.reload()
    }
  }
}
