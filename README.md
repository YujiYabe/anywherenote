# anywhereNote
  
- evernoteみたいなオンラインメモアプリ(見た目は全然似てませんが)  
- golang(echo,gorm)＋sqliteで作成  
- Electronのように内部でwebサーバをたててローカルのDBファイルにアクセス  
- 保存先はローカルPC内に無制限に追加可能
- ローカルのオンラインストレージサービス(dropbox等)に保存することで、そこにDBファイルを作成
- DBファイル自体のインターネット同期はオンラインストレージサービス(dropbox等)に丸投げ
- 1保存先=1DB=1ノート、ノートの中の1行=1ページ 


## 利用イメージ  
![利用イメージ](https://github.com/YujiYabe/anywhereNote/blob/garage/imageuse.jpg "")




## Screenshot  
![Screenshot](https://github.com/YujiYabe/anywhereNote/blob/garage/explain2.gif "")

---

[win64_v0.1]: https://drive.google.com/open?id=1W9S-JLfF8dgkO3fbLGOGDYkReJTm-lBb "Windows64_v0.1"
[win32_v0.1]: https://drive.google.com/open?id=1UqiawXaHZhSfxU5clmMt7JtBUG2pyzYk "Windows32_v0.1"
[lnx64_v0.1]: https://drive.google.com/open?id=1gLXapKzuW9U195F_C_DquuKgu1tFDXrc "Linux64_v0.1"
[lnx32_v0.1]: https://drive.google.com/open?id=19wQxlKyzaEFViVKLj9ID4J2DfRoMTkdb "Linux32_v0.1"

[win64_v0.2]: https://drive.google.com/open?id=141cNdQlNrW4H0lFWu_ib_4w8Vc4zVNm7 "Windows64_v0.2"
[win32_v0.2]: https://drive.google.com/open?id=11ogdDpNSyp7omn3r4GCp3y03UhUR5PPS "Windows32_v0.2"
[lnx64_v0.2]: https://drive.google.com/open?id=1HtcChZZ4CFFaoBB1VLgzYBXP1r5NUgKY "Linux64_v0.2"
[lnx32_v0.2]: https://drive.google.com/open?id=1F7EiJSrp2igFuBLMVmEKmRYnxmh8Cp0g "Linux32_v0.2"

[win64_v0.3]: https://drive.google.com/open?id=1tzOQHrnWySrIdpddK_fn4CqPwLj_4H42 "Windows64_v0.3"
[win32_v0.3]: https://drive.google.com/open?id=1_kgJwcRH6INMmoEcE132o7KIZI31dXvF "Windows32_v0.3"
[lnx64_v0.3]: https://drive.google.com/open?id=1CUXeo3Bqpv4dZ2bPPdICKCsTA7uATE9U "Linux64_v0.3"
[lnx32_v0.3]: https://drive.google.com/open?id=1wQQMQPkS-JnH3JRVvB4aV5jyhVMeZO8Z "Linux32_v0.3"


## ダウンロード

|      | Windows 64bit | Windows 32bit | Linux 64bit  | Linux 32bit  | Macintosh  |
|:-----|:--------------|:--------------|:-------------|:-------------|:-----------|
| v0.1 | [win64_v0.1]  | [win32_v0.1]  | [lnx64_v0.1] | [lnx32_v0.1] | 準備中     |
| v0.2 | [win64_v0.2]  | [win32_v0.2]  | [lnx64_v0.2] | [lnx32_v0.2] | 準備中     |
| v0.3 | [win64_v0.3]  | [win32_v0.3]  | [lnx64_v0.3] | [lnx32_v0.3] | 準備中     |
   
  


  

## evernote(フリー版)にはないメリット
- DBファイルを無限に追加できる。→容量制限はオンラインストレージサービス先に依存
- インストール台数無制限
- 1ソースでWindows・Linux・Mac向けのアプリにコンパイル可能(開発者視点)
  


## 今後のタスク
- [x] ファイルアップロード・ダウンロード機能追加  
- [ ] FireFox対応  
- [ ] Android、iOS対応（願望)  
- [ ] テストコード追加  
- [ ] react.js導入 
  
