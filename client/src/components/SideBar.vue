<template>
    <div class="sidebar">
        <h2>Bots List</h2>
        <ul class="bot-list">
            <li v-for="(value, index) in sessions" :key="index"
                :class="{current: currentSessionIndex === index}"
                @click="selectSession(index)">
                <div class="bot-avatar">
                    <img :src="value.avatar" alt="bot avatar"/>
                    <h6 style="white-space: nowrap;">{{ value.name }}</h6>
                </div>
                <div class="bot-info">
                    <p class="last-message">{{ value.messages[value.messages.length - 1].content }}</p>
                    <p class="chat-time">{{ value.messages[value.messages.length - 1].time }}</p>
                </div>
            </li>
        </ul>
    </div>
</template>

<script>
export default {
    name: 'SideBar',
    computed: {
        sessions () {
            return this.$store.state.sessions
        },
        currentSessionIndex () {
            return this.$store.state.currentSessionIndex
        }
    },
    methods: {
      selectSession (index) {
            this.$store.commit('SET_CURRENT_SESSION', index)
        },
    },
    // before the chat, use `sessions.json` to load bot_name and init_prompt
    created () {
        this.$store.dispatch('fetchInit')
    }
}
</script>

<style scoped>
.sidebar {
    position: absolute;
    width: 20%;
    background-color: skyblue;
    top: 0;
    left: 0;
    height: 100%;
    overflow-y: auto;
}

.bot-list {
    list-style: none;
    margin: 0;
    padding: 0;
    height: 100%; /* adjust this value to fit your layout */
    overflow-y: auto;
    background-color: lightblue;
}

.sidebar h2 {
    margin-top: 0;
}

.sidebar li {
    cursor: pointer;
    padding: 10px;
    border-radius: 10px;
}

.sidebar li.current {
    background-color: #e6e6e6;
}

.bot-avatar {
    width: 50px;
    height: 50px;
    display: flex;
    flex-direction: row;
    align-items: center;
}

.bot-avatar img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    float: left;
    border-radius: 100%;
    margin-right: 10px;
}

.bot-info {
    flex: 1;
}

.bot-info h3 {
    font-size: 16px;
    font-weight: bold;
    margin: 0;
}

.last-message {
    font-size: 14px;
    margin: 5px 0;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}

.chat-time {
    font-size: 12px;
    color: #888;
    margin: 0;
}

</style>
