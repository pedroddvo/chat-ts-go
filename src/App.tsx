import React, { useEffect, useState } from 'react';
import ChatBox from './component/ChatBox';
import MessageBox from './component/MessageBox';

const ws = new WebSocket("ws://localhost:8080/ws");

const App = () => {
  const [messages, setMessages] = useState<string[]>([]);

  useEffect(() => {
    ws.addEventListener('open', event => {
      console.log("Connection to websocket established.");
    });

    ws.addEventListener('message', m => {
      const msgs: string[] = JSON.parse(m.data).map(atob);
      setMessages(msgs)
    })
  }, [])

  return (
    <div className='w-full h-full absolute flex flex-col justify-between'>
      <MessageBox messages={messages} />
      <ChatBox ws={ws} />
    </div>
  );
}

export default App;
