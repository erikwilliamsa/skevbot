# SkevBot


## Environment Variables


###

The follow are required for the application to run:

- CLIENT_ID
- CLIENT_SECRET

Get tokens from [The dev console](https://dev.twitch.tv/console)

1. Click register Application
2. Ignore callbacks
3. Click generate secret

**Note** This secret will not be visible when you revisit, so you will need to generate a new one if you for get yours!

### Chat Environment Variables


- CHANNEL -  This is the channel you'll connect to in chat to monitor.

The following is only needed if you want to connect to your chat as a non anonymous user.  

- OAUTH_CHAT_TOKEN
- TWITCH_USERNAME


Get your `oauth token` from the [twitch generator](https://twitchapps.com/tmi/).