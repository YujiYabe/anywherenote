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




## Screenshots  
![Screenshots](https://github.com/YujiYabe/anywhereNote/blob/garage/explain2.gif "")

---

[win_64]: https://drive.google.com/open?id=1W9S-JLfF8dgkO3fbLGOGDYkReJTm-lBb "Windows 64"
[win_32]: https://drive.google.com/open?id=1UqiawXaHZhSfxU5clmMt7JtBUG2pyzYk "Windows 32"
[lnx_64]: https://drive.google.com/open?id=1gLXapKzuW9U195F_C_DquuKgu1tFDXrc "Linux 64"
[lnx_32]: https://drive.google.com/open?id=19wQxlKyzaEFViVKLj9ID4J2DfRoMTkdb "Linux 32"


## ダウンロード
|       | Windows  | Linux    | Macintosh  |
|:------|:---------|:---------|:-----------|
| 64bit | [win_64] | [lnx_64] | 準備中|
| 32bit | [win_32] | [lnx_32] | 準備中      |

 



  

## evernote(フリー版)にはないメリット
- DBファイルを無限に追加できる。→容量制限はオンラインストレージサービス先に依存
- インストール台数無制限
- 1ソースでWindows・Linux・Mac向けのアプリにコンパイル可能(開発者視点)
  


## 今後のタスク
- [ ] テストコード追加  
- [ ] ファイルアップロード機能追加  
- [ ] react.js導入 
- [ ] FireFox対応  
- [ ] Android、iOS対応（願望)  
  
