# anywhereNote
  
- evernoteみたいなオンラインメモアプリ  
- golang(echo,gorm)＋sqliteで作成  
- Electronのように内部でwebサーバをたててローカルのDBファイルにアクセス  
- 保存先を追加(dropbox等)することで、その保存先にDBファイルを作成
- DBファイル自体のインターネット同期はオンラインストレージサービス(dropbox等)に丸投げ
  

## 利用イメージ  
**ブラウザ⇔ローカルwebサーバー⇔ローカルDBファイル⇔オンラインストレージサービス**  
  
## Screenshots  
![Screenshots](https://github.com/YujiYabe/anywhereNote/blob/garage/explain.gif "")


|       | Windows | Linux | Machintosh |
|:------|:--------|:------|:-----------|
| 32bit | [Windows 64](https://drive.google.com/open?id=1W9S-JLfF8dgkO3fbLGOGDYkReJTm-lBb "Windows 64") | [Linux 32](https://drive.google.com/open?id=19wQxlKyzaEFViVKLj9ID4J2DfRoMTkdb "Linux 32") | 準備中      |
| 64bit | [Windows 32](https://drive.google.com/open?id=1UqiawXaHZhSfxU5clmMt7JtBUG2pyzYk "Windows 32") | [Linux 64](https://drive.google.com/open?id=1gLXapKzuW9U195F_C_DquuKgu1tFDXrc "Linux 64") | 準備中      |

  
  
## ダウンロード




  

## evernote(フリー版)にはないメリット
- DBファイルを無限に追加できる。→容量制限はオンラインストレージサービス先に依存
- 台数無制限
- 1ソースでWindows・Linux・Mac向けのアプリにコンパイル可能(開発者視点)
- 将来的にAndroid、iOS向けアプリに対応（願望）
  


## 今後のタスク
- テストコード追加
- ファイルアップロード機能追加
- FireFox対応  
- スマホ対応
  
