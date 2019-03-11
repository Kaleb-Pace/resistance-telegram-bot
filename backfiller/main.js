var fs = require("fs");
var content = JSON.parse(fs.readFileSync("result.json"));

let allChats = content.chats.list;

const chatToCleanID = 9771743476;

function CleanChat(chat) {
  for (var m = 0; m < chat.messages.length; m++) {
    if (Array.isArray(chat.messages[m].text)) {
      let newText = "";

      for (var t = 0; t < chat.messages[m].text.length; t++) {
        if (typeof chat.messages[m].text[t] === "string") {
          newText += chat.messages[m].text[t];
        } else if (typeof chat.messages[m].text[t] === "object") {
          newText += chat.messages[m].text[t].text;
        } else {
          console.error("Unkown Type", chat.messages[m].text[t]);
        }
      }

      chat.messages[m].text = newText;
    }
  }
}

for (var i = 0; i < allChats.length; i++) {
  if (allChats[i].id == chatToCleanID) {
    console.log(allChats[i].name);
    CleanChat(allChats[i]);
  }
}

fs.writeFile(
  "result-cleaned.json",
  JSON.stringify(content, null, 2),
  "utf8",
  function() {
    console.log("Finished");
  }
);
