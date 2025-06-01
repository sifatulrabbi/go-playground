import React, { useEffect, useState } from "react";

export const App: React.FC = () => {
  const [updatedAfter, setUpdatedAfter] = useState("");
  const [responseText, setResponseText] = useState("");

  let lastTs = 0;
  let count = 0;
  const loop = (ts: number) => {
    if (ts - lastTs >= 1000) {
      setUpdatedAfter(`Frame count: ${count}`);
      lastTs = ts;
      count = 0;
    }
    count++;

    window.requestAnimationFrame(loop);
  };

  const handleInitiateWs = async () => {
    const ws = new WebSocket("ws://localhost:8000/ws");
    ws.onopen = function (ev) {
      console.log("New connection established:", ev);
      ws.send(JSON.stringify({ message: "Hello world!" }));
    };
    ws.onmessage = function (ev) {
      setResponseText(ev.data);
    };
    ws.onclose = function (ev) {
      console.log("Connection closed.", ev);
    };
  };

  useEffect(() => {
    const animId = window.requestAnimationFrame(loop);
    return () => {
      window.cancelAnimationFrame(animId);
    };
  }, []);

  return (
    <div className="w-full flex flex-col min-h-screen min-w-screen max-w-screen p-6 gap-6">
      <h1 className="text-3xl font-bold">Hello World</h1>
      <p>{updatedAfter}</p>
      {responseText && <p>{responseText}</p>}
      <button
        onClick={handleInitiateWs}
        className="w-max rounded-md border-2 border-blue-500 bg-blue-500 text-white py-1 px-2 cursor-pointer hover:bg-blue-700 hover:border-blue-700 transition-colors"
      >
        Initiate WebSocket
      </button>
    </div>
  );
};
