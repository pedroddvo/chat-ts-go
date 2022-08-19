import React, { useState } from 'react';

interface ChatBoxProps {
    ws: WebSocket
}

const ChatBox = ({ ws }: ChatBoxProps) => {
    const [message, setMessage] = useState("")

    return (
        <div className='fixed flex flex-col items-center h-1/4 rounded-lg border-2 m-8 p-4 gap-4 inset-x-0 bottom-0 bg-white shadow-lg'>
            <div className='flex w-full justify-between'>
                <p className='font-bold text-xl'>Message</p>
                <button className='rounded-full bg-green-500 px-4 py-1 text-green-100'
                    onClick={() => ws.send(message)}>Send Message</button>
            </div>
            <textarea className='w-full h-full'
                onChange={v => setMessage(v.target.value)} />
        </div>
    );
}

export default ChatBox;