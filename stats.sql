# Get all messages
select * FROM db.messages;

# This counts how many messages each person sends
SELECT COUNT(db.messages.MessageID), FromUserName 
	FROM db.messages
	WHERE ChatID = -1001181808884
	GROUP BY FromUserName;
 
# This counts how many messages each person forwards
SELECT COUNT(db.messages.ForwardedFromUserID), FromUserName 
	FROM db.messages
	GROUP BY FromUserName;
 
# Count all messages from a specific chat
SELECT COUNT(*) FROM db.messages WHERE ChatID = -1001181808884;
