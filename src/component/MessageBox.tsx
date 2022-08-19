import React from 'react';

interface MessageboxProps {
    messages: string[]
}

const MessageBox = ({ messages }: MessageboxProps) => {

    return (
        <div className='flex flex-col'>
            {messages.map(
                (s, i) => 
                    <div key={`post-${i}`} className={`p-2 bg-indigo-${(i % 2 + 2) * 100}`}>
                        <p className='font-bold'>Anonymous</p>
                        <p className={`p-2 `}>{s}</p>
                    </div>
            )}
        </div>
    )
}

export default MessageBox;