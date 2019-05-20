# Get all messages
select * FROM db.messages;

# This counts how many messages each person sends
SELECT COUNT(db.messages.MessageID) as c, FromID, FromUserName
	FROM db.messages
	WHERE ChatID = -1001181808884 AND Text is not Null
	GROUP BY FromUserName, FromID
    order by c desc;
  
# This counts how many messages each person forwards
SELECT COUNT(db.messages.PhotoFileID), FromUserName 
	FROM db.messages
	GROUP BY FromUserName;
    
SELECT * FROM db.messages WHERE ChatID != -1001181808884 AND ChatID != 106468411;

SELECT * FROM db.messages WHERE ChatID = -1001181808884;

# Get number of messages sent per chat.
SELECT COUNT(db.messages.ChatID), ChatID
	FROM db.messages
	GROUP BY ChatID;

SELECT ReplyToMessageID, COUNT(db.messages.MessageID) AS num
          FROM db.messages GROUP BY ReplyToMessageID
          ORDER BY num DESC;
    
SELECT db.messages.Text FROM db.messages WHERE PhotoFileID is not NULL;
 
# Count all messages from a specific chat
SELECT COUNT(*) FROM db.messages WHERE ChatID = -256262125;

SELECT COUNT(db.messages.MessageID), FromUserName FROM db.messages WHERE db.messages.Text LIKE "holy shit" GROUP BY FromUserName;

SELECT * FROM db.messages WHERE ChatID = -256262125;
