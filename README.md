# insert-google-calendar

 - [wataruの作ってたやつ](https://github.com/flat35hd99/play-oauth)をベースにしてるよ
    - 認証コードが自動で取得されて便利
 - 認証のときブラウザが自動で開いたり、calendarに予定が作れるように改良を加えたよ
 - 使うときは[ここ](https://console.cloud.google.com/apis/credentials)からtestの方をダウンロードして、credentials.jsonとして保存してね
 - events.csvっていう名前のcsvファイルを自動で読み込んでそれに基づいてイベントを作るようにしてあるよ
 - csvはsummary,start_datetime,end_datetime,timezone,attendeesの形式で入れるようになってるよ
    - attendeeは;で分けて複数入れられるよ
    - 会議室を追加したいときは、attendeeに会議室のcalenderidを入れればいけるよ
    - 複数行入れれば、複数イベントが同時に入るよ
 - csvを読み込むメソッドの内容とかcsvの形式とかに関しては改造してよしなに変えてほしい