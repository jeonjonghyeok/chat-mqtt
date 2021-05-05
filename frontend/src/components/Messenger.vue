<template>
  <MenuBar 
  :roomName="props.roomName" 
  :roomID="props.roomID"
  :token="props.token"
  @disconnect="disconnect" 
  @contract="getContract"
  />
  <div class="messages"
  v-if="!getOpen"
  onscroll="chat_on_scroll()"
  >
    <Message
      v-for="(message, idx) in messages"
      :message="message"
      :mine="myID === message.senderID"
      :displaySender="(myID !== message.senderID) && (messages[idx-1] ?? {sender: null}).senderID !== message.senderID"
      :key="idx"
    />
  </div>
  <div class="message-input-bar"
  v-if="!getOpen"
  >
    <MessageInput @send="onSend" />
  </div>
<div 
  class="get-contract"
    v-else
>
<h1 class="contract-title"
> 계약서 </h1>
    <Contracts
      v-for="(contract, idx) in contracts"
      :contract="contract"
      :key="idx"
    />
    <Contract
      v-for="(message, idx) in contract"
      :message="message"
      :mine="myID === message.senderID"
      :displaySender="(contract[idx-1] ?? {sender: null}).senderID !== message.senderID"
      :key="idx"
    />
  </div>
</template>

<script>
import { ref } from "vue";
import Message from "./Message.vue";
import MessageInput from "./MessageInput.vue";
import MenuBar from "./MenuBar.vue";
import Contract from "./Contract.vue";
import Contracts from "./Contracts.vue";

export default {
  components: { Message, MessageInput, MenuBar, Contracts, Contract },
  name: "Messenger",
  emits: ['disconnect','contract'],
  props: {
    token: String,
    roomID: Number,
    roomName: String
  },
  setup(props, {emit}) {
    const myID = ref('');
    const socket = new WebSocket(`ws://110.165.17.149:5000/room/${props.roomID}?token=${encodeURIComponent(props.token)}`);
    const messages = ref([]);
    const contract = ref([]);
    const contracts = ref([]);
    const getOpen = ref(false)
    const name = ref("")
    
    const getContracts = async () => {
      const resp = await fetch(`http://110.165.17.149:5000/getContracts?token=${encodeURIComponent(props.token)}`, {
          method: 'POST',
          body: JSON.stringify({
            'roomID': props.roomID
          })
        })
        const parsed = await resp.json()
        getOpen.value=true
        name.value = props.roomName
        contracts.value=parsed
        getContract
    }

    const getContract = async () => {
      const resp = await fetch(`http://110.165.17.149:5000/getContract?token=${encodeURIComponent(props.token)}`, {
          method: 'POST',
          body: JSON.stringify({
            'roomID': props.roomID
          })
        })
        const parsed = await resp.json()
        contract.value = parsed
        getOpen.value=true
        name.value = props.roomName
    }


    const parseJwt = (token) => {
      try {
        return JSON.parse(atob(token.split('.')[1]));
      } catch (e) {
        return null;
      }
    };

    const parsedToken = parseJwt(props.token)
    myID.value = parsedToken.uid
    const onSend = message =>
      socket.send(
        JSON.stringify({
          text: message,
        })
      );
    socket.onmessage = messageEvent =>
      messages.value.push(JSON.parse(messageEvent.data));
    return {
      myID,
      messages,
      onSend,
      props,
      disconnect: () => emit('disconnect'),
      getContract,
      getOpen,
      contract,
      getContracts,
      contracts
    };
  }
};
</script>

<style scoped>
.get-contract{
  width: 500px;
  height: 100%;
  margin: auto;

  display: flex;
  flex-direction: column;
  justify-content: space-between;
  align-items: stretch;

  flex-grow: 1;
  overflow: auto;
  padding: 10px;
  display: flex;
  flex-direction: column;
  align-items: stretch;

  padding-top: 70px;
  padding-bottom: 70px;
}
.get-contracts{
  width: 500px;
  height: 100%;
  margin: auto;

  display: flex;
  flex-direction: column;
  justify-content: space-between;
  align-items: stretch;

  flex-grow: 1;
  overflow: auto;
  padding: 10px;
  display: flex;
  flex-direction: column;
  align-items: stretch;

  padding-top: 70px;
  padding-bottom: 70px;
}
.messages {
  width: 500px;
  height: 100%;
  margin: auto;

  display: flex;
  flex-direction: column;
  justify-content: space-between;
  align-items: stretch;

  flex-grow: 1;
  overflow: auto;
  padding: 10px;
  display: flex;
  flex-direction: column;
  align-items: stretch;

  padding-top: 70px;
  padding-bottom: 70px;
}
.contract-title {
  font-family: monospace;
  font-size: 30px;
  text-align: center;
}

.message-input-bar {
  height: 60px;
  width: 100%;
  position: fixed;
  bottom: 0;
  background-color: #bbc9e0;
}
</style>