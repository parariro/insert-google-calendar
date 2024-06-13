# insert-google-calendar

 - [wataruの作ってたやつ](https://github.com/flat35hd99/play-oauth)をベースにしてるよ
    - 認証コードが自動で取得されて便利
 - 認証のときブラウザが自動で開いたり、calendarに予定が作れるように改良を加えたよ
 - クレデンシャルが必要だから、[ここ](https://console.cloud.google.com/apis/credentials)からtestの方をダウンロードしてcredentials.jsonとして保存して
 - events.csvっていう名前のcsvファイルを読み込んでイベントを作成するようにしてあるよ
 - csvはsummary,start_datetime,end_datetime,timezone,attendeesの形式になってるよ
    - attendeeは;で分けて複数入れられるよ
    - 会議室を追加したいときは、attendeeに会議室のcalenderidを入れればいけるよ
        - 会議室のcalendaridは、他のカレンダーに会議室を追加して点々のとこから設定に行けば見られるよ
    - 複数行書いたら複数イベントが同時に入るよ
 - csvを読み込むメソッドの内容とかcsvの形式とかに関しては改造してよしなに変えてほしい