# AlertInfo
A plugin with ChatGPT and Google integration for biker_info system, that allows web searching and AI content processing

AlertInfo (AI) is another sms service of my, that allows for ordering any information from the internet, in a chosen period of time.

You can try the service yourself under the telephone number (0048) 579 390 848

In the beginning of your message, you should include notification frequency. 
It can be chosen from the following:
'alert1d' - to receive the desired information every day at 6pm
'alerttr' - to receive the message within the workweek 
'alert7d' - to receive the message every Friday at 6pm
'alert30d' - to receive it only on 28th months day

Exemplary use-cases:
'Alerttr Tesla stock price'
or
'Alert30d next full moon date'

The challenge:
Unlike it is, using the web application, openAI GPT API does not allow for searching on the Internet.
The walk-around:
To bypass this incapability, I propose the following processes of gaining a similar result, by combining Google search API with text analysis provided by openAI chatGPT API:
![image](https://github.com/user-attachments/assets/f2eea8de-f393-49b6-8b9b-776144c26421)


Contracts examples for each API:

Google Search:
curl -X GET "https://www.googleapis.com/customsearch/v1?q=your_query&key=GOOGLE_SEARCH_API_KEY&cx=YOUR_SEARCH_ENGINE_ID"

openAI:
curl "https://api.openai.com/v1/chat/completions"     -H "Content-Type: application/json"     -H "Authorization: Bearer $OPENAI_API_KEY"     -d '{
        "model": "gpt-4o-mini",
        "messages": [
            {
                "role": "system",
                "content": "You are a helpful assistant."
            },
            {
                "role": "user",
                "content": "as an answer, provide me ONLY a link to a web service, where i can find an answer to my question: Tesla stock price"
            }
        ]
    }'


