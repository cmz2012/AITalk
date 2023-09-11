<template>
    <div>
        <Header></Header>
        <section ref="messagesRef" class="chatList">
            <ul>
                <div v-for="(message, index) in bots[currentBotIndex].messages" :key="index">
                    <li v-if="message.sender ==='user' " class="chat-user">
                        <div class="chat-img"><img src="assets/user.png"></div>
                        <div class="time"><cite>{{ message.sender }}<i>{{ message.time }}</i></cite></div>
                        <div class="chat-content" style="white-space: pre-wrap;" v-text="message.content" v-on:click="playAudio(message.audio)"></div>
                    </li>
                    <li v-if="message.sender !== 'user'" class="chat-others">
                        <div class="chat-img"><img :src="bots[currentBotIndex].avatar"></div>
                        <div class="time"><cite>{{ message.sender }}<i>{{ message.time }}</i></cite></div>
                        <div class="chat-content" style="white-space: pre-wrap;" v-text="message.content" v-on:click="playAudio(message.audio)"></div>
                    </li>
                </div>
            </ul>
        </section>

        <section class="foot">
            <mt-field id="txtContent" v-model="content" class="con" placeholder="Please click start to record audio.">
            </mt-field>
            <span class="btn btn-send" v-on:click="recordMsg">
                {{state}}
            </span>
        </section>
    </div>
</template>

<script>
import util from '../common/util'
// import { Toast } from 'mint-ui'
// import axios from 'axios'
import Header from './Header.vue'
import Recorder from 'js-audio-recorder';

export default {
    // eslint-disable-next-line vue/multi-word-component-names
    name: 'ChatList',
    components: {
        Header
    },
    data () {
        return {
            content: '',
            records: [], // chat record
            ws: null,
            recorder: null,
            audio: null,
            state: 'start'
        }
    },
    created () {
        this.ws = new WebSocket('ws://localhost:8888/chat')
        this.ws.onmessage = this.onMessage
        console.log('ws init')
        this.recorder = new Recorder({
            sampleBits: 16, // 采样位数，支持 8 或 16，默认是16
            sampleRate: 16000, // 采样率，支持 11025、16000、22050、24000、44100、48000，根据浏览器默认值，我的chrome是48000
            numChannels: 1, // 声道，支持 1 或 2， 默认是1
            // compiling: false,(0.x版本中生效,1.x增加中)  // 是否边录边转换，默认是false
        })
        console.log("recorder init")

        this.audio = document.createElement("audio"); //创建标签
        this.audio.style.display = "none";
        console.log('audio init')
    },
    computed: {
        currentBotIndex () {
            return this.$store.state.currentBotIndex
        },
        bots () {
            return this.$store.state.bots
        }
    },
    methods: {
        // 播放音频
        playAudio(url) {
          this.audio.src = url
            this.audio.play()
        },
        //开始录音
        startRecordAudio() {
            Recorder.getPermission().then(
                () => {
                    console.log("开始录音");
                    this.recorder.start(); // 开始录音
                },
                (error) => {
                    console.log(`${error.name} : ${error.message}`);
                }
            );
        },
        //停止录音
        stopRecordAudio() {
            console.log("停止录音");
            this.recorder.stop();
        },
        // 接受websocket消息
        onMessage: function (e) {
            const newbolb = new Blob([e.data], { type: 'audio/wav' })
            const audioSrc = URL.createObjectURL(newbolb);
            this.$store.commit('SEND_MESSAGE', {
                content: 'audio_text',
                sender: 'bot',
                receiver: this.bots[this.currentBotIndex].name,
                time: util.formatDate.format(new Date(), 'yyyy-MM-dd hh:mm:ss'),
                audio: audioSrc
            })
            this.playAudio(audioSrc)

            this.scrollToBottom()
            this.focusTxtContent()
        },
        // 录制消息
        recordMsg: function () {
            if (this.state === 'start') {
                // 开始录音
                this.startRecordAudio()
                this.state = 'done'
            } else {
                // 完成录音
                this.stopRecordAudio()
                var blob = this.recorder.getWAVBlob()
                const audioSrc = URL.createObjectURL(blob);

                // this.playAudio(audioSrc)

                this.$store.commit('SEND_MESSAGE', {
                    content: 'transcribing...',
                    sender: 'user',
                    receiver: this.bots[this.currentBotIndex].name,
                    time: util.formatDate.format(new Date(), 'yyyy-MM-dd hh:mm:ss'),
                    audio: audioSrc
                })
                this.state = 'start'

                // Make the dialog box go to the bottom
                // Maybe it is Duplicate
                this.$nextTick(() => {
                    this.$refs.messagesRef.scrollTop = this.$refs.messagesRef.scrollHeight
                })

                // write to server
                this.ws.send(blob)
            }
        },

        // focus the input box
        focusTxtContent: function () {
            document.querySelector('#txtContent input').focus()
        },

        // Scroll bar scrolls to the bottom
        scrollToBottom: function () {
            setTimeout(function () {
                const tempChatList = document.getElementsByClassName('chatList')[0]
                tempChatList.scrollTop = tempChatList.scrollHeight
            }, 100)
        },
    },
    mounted: function () {
        this.scrollToBottom()
        this.focusTxtContent()
    }
}
</script>

<style lang="css" scoped>
.chatList {
    position: absolute;
    top: 48px;
    bottom: 48px;
    left: 17%;
    right: 0;
    overflow-y: scroll;
    overflow-x: hidden;
    padding: 20px;
}

.chatList ul {
    min-height: 300px;
}

.chatList ul .chat-user {
    text-align: right;
    padding-left: 0;
    padding-right: 60px;
}


.chatList ul .chat-others {
    text-align: left;
    padding-left: 60px;
    padding-right: 0;
}

.chatList ul li {
    position: relative;
    margin-bottom: 10px;
    padding-left: 60px;
    min-height: 68px;
    /*去掉li的默认小圆点*/
    list-style-type: none;
}

.chat-user .chat-img {
    left: auto;
    right: 3px;
}

.chat-others .chat-img {
    left: 3px;
    right: auto;
}

.chat-img {
    position: absolute;
    left: 3px;
}

.chat-content,
.chat-img {
    display: inline-block;
    vertical-align: top;
    /*font-size: 14px;*/
}

.chat-img img {
    width: 40px;
    height: 40px;
    border-radius: 100%;
}

.time {
    width: 100%;
}

cite {
    left: auto;
    right: 60px;
    text-align: right;
}

cite {
    font-style: normal;
    line-height: 24px;
    font-size: 12px;
    white-space: nowrap;
    color: #999;
    text-align: left;
}

cite i {
    font-style: normal;
    padding-left: 5px;
    padding-right: 5px;
    font-size: 12px;
}

.chat-user .chat-content {
    margin-left: 0;
    text-align: left;
    background-color: #33DF83;
    color: #fff;
}

.chat-others .chat-content {
    margin-left: 0;
    text-align: left;
    background-color: #33DF83;
    color: #fff;
}

.chat-content {
    position: relative;
    line-height: 22px;
    /*margin-top: 25px;*/
    padding: 10px 15px;
    background-color: #eee;
    border-radius: 3px;
    color: #333;
    word-break: break-all;
    max-width: 462px \9;
}

.chat-content,
.chat-img {
    display: inline-block;
    vertical-align: top;
    font-size: 14px;
}

.chat-content img {
    max-width: 100%;
    vertical-align: middle;
}

.chat-img {
    position: absolute;
    left: 3px;
}

.chat-content:after {
    content: '';
    position: absolute;
    left: -10px;
    top: 15px;
    width: 0;
    height: 0;
    border-style: solid dashed dashed;
    border-color: #eee transparent transparent;
    overflow: hidden;
    border-width: 10px;
}

.chat-user .chat-content:after {
    left: auto;
    right: -10px;
    border-top-color: #33DF83;
}

.chat-others .chat-content:after {
    left: -10px;
    right: auto;
    border-top-color: #33DF83;
}

.foot {
    width: 80%;
    min-height: 48px;
    position: fixed;
    bottom: 0px;
    right: 0px;
    background-color: #F8F8F8;
}

.foot .con {
    position: absolute;
    left: 0px;
    right: 0px;
}

.btn {
    display: inline-block;
    vertical-align: top;
    font-size: 30px;
    line-height: 48px;
    margin-left: 5px;
    padding: 0 6px;
    background-color: #33DF83;
    color: #fff;
    border-radius: 3px;
}

.btn-send {
    position: absolute;
    right: 0px;
}
</style>
