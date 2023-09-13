// Import and use Vuex
import Vuex from 'vuex'
import Vue from 'vue'
// import util from '@/common/util'
import axios from 'axios'

Vue.use(Vuex)

// actions in the corresponding component
const actions = {
    // before the chat, use `bots.json` to load bot_name and init_prompt
    async fetchInit (context) {
        const session_ids = await context.dispatch('fetchSession')
        // console.log("session_ids: ", session_ids)
        for (const s of session_ids) {
            const messages = await context.dispatch('fetchSessionMsg', s)
            context.state.sessions[s] = {
                name: 'bot',
                avatar: 'assets/' + Math.floor(Math.random() * 10) + '.png',
                messages: []
            }
            // console.log("session: ", context.state.sessions)
            // console.log("session_message: ", messages)
            for (const message of messages) {
                context.commit('SEND_MESSAGE', {
                    content: message['data'],
                    sender: message['user_id'],
                    time: message['create_time'],
                    session: s,
                    audio: message['audio_key'],
                })
            }
        }
        context.commit('SET_CURRENT_SESSION', session_ids[0])
    },

    async fetchSession(context) {
        // 从服务端拉去user下的session
        const param = {
            'user_id': context.state.user_id
        }
        const path = `http://${window.location.hostname}:8888/session`;
        var sessionIDs = []
        await axios.get(path, {params: param})
            .then((res) => {
                // console.log(res.data['sessions'])
                sessionIDs.push(...res.data['sessions'])
            })
            .catch((error) => {
                console.log(error)
                return []
            })
        return sessionIDs
    },

    async fetchSessionMsg(context, session_id) {
        const paramMessage = {
            'user_id': context.state.user_id,
            'session_id': session_id
        }
        var msg = []
        const path = `http://${window.location.hostname}:8888/message`
        await axios.get(path, {params: paramMessage}).then((res) => {
            // console.log(res.data['messages'])
            msg.push(...res.data['messages'])
        }).catch((error) => {
            console.log(error)
            return []
        })
        return msg
    }
}

// manipulating data in `state`
const mutations = {
    SEND_MESSAGE (state, payload) {
        // console.log('SEND_MESSAGE in mutations is called.')
        // bot to user
        state.sessions[payload.session].messages.push(payload)
    },

    SET_CURRENT_SESSION (state, value) {
        // console.log('SET_CURRENT_SESSION in mutation is called.')
        state.currentSessionIndex = value
    },

    GET_ALL_SESSION (state) {
        try {
            const session_ids = this.commit('FETCH_SESSION')
            console.log("session_ids: ", session_ids)
            for (const s of session_ids) {
                const messages = this.commit("FETCH_SESSION_MSG", s)
                state.sessions[s] = {
                    name: 'bot',
                    avatar: 'assets/' + Math.floor(Math.random() * 10) + '.png',
                    messages: []
                }
                console.log("session message: ", messages)
                for (const message of messages) {
                    this.commit('SEND_MESSAGE', {
                        content: message['data'],
                        sender: message['user_id'],
                        time: message['create_time'],
                        session: s,
                        audio: message['audio_key'],
                    })
                }
            }
            console.log("session: ", state.sessions)
        } catch (error) {
            console.error(error)
        }
    },

    // before the chat, use `sessions.json` to load bot_name and init_prompt
    INIT_SESSIONS: function (state) {
        console.log("state.sessions: ", state.sessions)
    }
}

// Storing data
const state = {
    currentSessionIndex: 0,
    sessions: {},
    user_id: 1
}

// Create and export Store
export default new Vuex.Store({
    actions: actions,
    mutations: mutations,
    state: state
})
