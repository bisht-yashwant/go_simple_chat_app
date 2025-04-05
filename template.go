package main

import "html/template"

var html = template.Must(template.New("").Parse(`
{{define "join_page"}}
<!DOCTYPE html><html><head><meta charset="UTF-8"><title>Join Chat</title>
<style>
  body{display:flex;justify-content:center;align-items:center;height:100vh;margin:0;font-family:sans-serif;background:#f0f0f0}
  form{background:white;padding:20px;border-radius:8px;box-shadow:0 0 10px rgba(0,0,0,0.1);width:280px}
  input,button{width:-webkit-fill-available;width: -moz-available;margin:8px 0;padding:10px;font-size:14px}
  button{background:#075e54;color:white;border:none;cursor:pointer;border-radius:4px}
</style>
<script>
  function saveName(){var u=document.getElementById('username').value;localStorage.setItem('chat_user',u);}
</script>
</head><body>
<form action="/join" method="POST" onsubmit="saveName()">
  <h3>Join Chat Room</h3>
  <input id="username" name="username" type="text" placeholder="Your Name" required>
  <input name="roomid" type="text" placeholder="Room ID" required>
  <button type="submit">Join</button>
</form>
</body></html>
{{end}}

{{define "chat_room"}}
<!DOCTYPE html><html><head><meta charset="UTF-8"><title>Room: {{.roomid}}</title>
<style>
  body{margin:0;font-family:'Segoe UI',sans-serif;background:#e5ddd5}
  .container{max-width:600px;margin:auto;height:100vh;display:flex;flex-direction:column;background:white}
  .header{background:#075e54;color:white;padding:16px;text-align:center}
  .msgs{flex:1;padding:16px;overflow-y:auto;display:flex;flex-direction:column}
  .msg{max-width:70%;margin:4px 0;padding:8px 12px;border-radius:8px}
  .self{background:#dcf8c6;align-self:flex-end}
  .other{background:white;border:1px solid #ddd;align-self:flex-start}
  .footer{display:flex;padding:12px;background:#f0f0f0}
  input[type=text]{flex:1;padding:10px;border:none;border-radius:20px;outline:none}
  button{margin-left:8px;padding:10px 16px;background:#075e54;color:white;border:none;border-radius:20px;cursor:pointer}
</style>
<script src="https://code.jquery.com/jquery-3.6.0.min.js"></script>
<script>
  var userID = localStorage.getItem('chat_user') || "{{.userid}}";
  $(function(){
    var src = new EventSource('/stream/{{.roomid}}');
    src.onmessage = function(e){
      var parts=e.data.split(': ');
      var usr=parts[0], txt=parts.slice(1).join(': ');
      var cls=(usr===userID)?'msg self':'msg other';
      $('.msgs').append('<div class="'+cls+'"><strong>'+usr+'</strong><br>'+txt+'</div>');
      $('.msgs').scrollTop($('.msgs')[0].scrollHeight);
    };
  });
  function sendMsg(){
    var txt=$('#m').val(); if(!txt)return;
    $.post('/room/{{.roomid}}',{user:userID,message:txt});
    $('#m').val('').focus();
  }
</script>
</head><body>
<div class="container">
  <div class="header">Chat Room: {{.roomid}}</div>
  <div class="msgs"></div>
  <div class="footer">
    <input id="m" type="text" placeholder="Type a message..." autocomplete="off">
    <button onclick="sendMsg()">Send</button>
  </div>
</div>
</body></html>
{{end}}
`))
