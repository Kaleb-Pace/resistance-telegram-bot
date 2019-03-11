# Get all messages
select * FROM db.messages;

# This counts how many messages each person sends
SELECT COUNT(db.messages.MessageID), FromID, FromUserName
	FROM db.messages
	WHERE ChatID = -1001181808884 AND Text is not Null
	GROUP BY FromUserName, FromID;
  
# This counts how many messages each person forwards
SELECT COUNT(db.messages.PhotoFileID), FromUserName 
	FROM db.messages
	GROUP BY FromUserName;
    
    
SELECT db.messages.Text FROM db.messages WHERE PhotoFileID is not NULL;
 
# Count all messages from a specific chat
SELECT COUNT(*) FROM db.messages WHERE ChatID = -1001181808884;

SELECT COUNT(db.messages.MessageID), FromUserName FROM db.messages WHERE db.messages.Text LIKE "holy shit" GROUP BY FromUserName;

