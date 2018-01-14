# anywhereNote
  
- evernoteみたいなオンライン同期メモアプリ  
- golang(echo,gorm)＋sqliteで作成  
- Electronのように内部でwebサーバをたててローカルのDBファイルにアクセス  
- 保存先を追加(dropbox等)することで、その保存先にDBファイルを作成
- DBファイル自体のインターネット同期はオンラインストレージサービス(dropbox等)に丸投げ
  

## 利用イメージ  
![利用イメージ](https://github.com/YujiYabe/anywhereNote/blob/garage/imageuse.jpg "")
※1保存先＝1DB(ノート)



## Screenshots  
![Screenshots](https://github.com/YujiYabe/anywhereNote/blob/garage/explain2.gif "")

---

[z][aaa]

## ダウンロード
|       | Windows | Linux | Macintosh |
|:------|:--------|:------|:----------|
| 64bit | [Windows 64](https://drive.google.com/open?id=1W9S-JLfF8dgkO3fbLGOGDYkReJTm-lBb "Windows 64") | [Linux 64](https://drive.google.com/open?id=1gLXapKzuW9U195F_C_DquuKgu1tFDXrc "Linux 64") | 準備中      |
| 32bit | [Windows 32](https://drive.google.com/open?id=1UqiawXaHZhSfxU5clmMt7JtBUG2pyzYk "Windows 32") | [Linux 32](https://drive.google.com/open?id=19wQxlKyzaEFViVKLj9ID4J2DfRoMTkdb "Linux 32") | 準備中      |

  
  




  

## evernote(フリー版)にはないメリット
- DBファイルを無限に追加できる。→容量制限はオンラインストレージサービス先に依存
- インストール台数無制限
- 1ソースでWindows・Linux・Mac向けのアプリにコンパイル可能(開発者視点)
  


## 今後のタスク
- テストコード追加
- ファイルアップロード機能追加
- React.js化 
- FireFox対応  
- Android、iOS対応（願望）
  
