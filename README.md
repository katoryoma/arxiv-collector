# arXiv論文収集・整理ツール 📚✨

arXiv APIを使用して論文を検索・収集し、自動的に整理するツールです♡

## 🎯 機能

### 1. 論文収集 (Go)
- キーワード検索で論文を取得
- 複数単語対応
- 最大100件まで取得可能
- CSV形式で保存

### 2. 論文整理 (Python)
- キーワード別にフォルダ分け
- 年度別にサブフォルダ作成
- 重複論文の自動除外
- 元のCSVをoriginalフォルダに保存

## 📁 プロジェクト構成

```
arxiv-collector/
├── main.go              # 論文収集メインプログラム
├── go.mod              # Go依存関係管理
├── box.py              # CSV整理スクリプト
├── show_csv.py         # CSV表示スクリプト
├── requirements.txt    # Python依存関係
├── README.md          # このファイル
└── paper_box/         # 整理済み論文フォルダ
    ├── original/      # 元のCSVファイル
    ├── キーワード1/
    │   ├── 2025/
    │   └── 2024/
    └── キーワード2/
        └── 2025/
```

## 🚀 セットアップ

### Go環境
```bash
# Go 1.21以上が必要
go version

# 依存関係は標準ライブラリのみ
```

### Python環境
```bash
# Python 3.8以上が必要
python --version

# 依存関係のインストール
pip install -r requirements.txt
```

## 💡 使い方

### 1. 論文を検索・収集
```bash
go run main.go
```

実行すると以下の順に入力を求められます：
1. **検索キーワード** - 例: `Machine Learning`, `face swap`
2. **取得件数** - 1〜100の範囲で指定
3. **続けるか** - `y`で次の検索、その他で終了

### 2. CSVを整理
```bash
python box.py
```

自動的に以下を実行：
- キーワード別フォルダ作成
- 年度別サブフォルダ作成
- 重複論文の除外
- 元CSVをoriginalフォルダに移動

### 3. CSVを表示（オプション）
```python
# show_csv.py を編集してパスを指定
python show_csv.py
```

## 📊 出力形式

### CSV列
| 列名 | 説明 |
|------|------|
| Keyword | 検索キーワード |
| Title | 論文タイトル |
| URL | arXiv URL |
| Summary | アブストラクト |
| Authors | 著者名（セミコロン区切り） |
| Published | 公開日時 |
| Publication | 出版情報（ジャーナル等） |

## 🔧 トラブルシューティング

### Goのコンパイルエラー
```bash
go mod tidy
go build main.go
```

### Pythonのインポートエラー
```bash
pip install --upgrade -r requirements.txt
```

### CSVが見つからない
- `go run main.go`を先に実行してCSVを生成してください
- カレントディレクトリを確認してください

## 📝 ワークフロー例

```bash
# 1. 論文を収集
go run main.go
# > 検索キーワード: Deep Learning
# > 取得件数: 50
# > 続ける: n

# 2. CSVを整理
python box.py
# > ✨ 完了しちゃった〜！

# 3. 結果を確認
ls paper_box/Deep_Learning/2025/
```

## 🎀 注意事項

- arXiv APIの利用制限に注意してください
- 大量リクエストは避けてください（1回あたり最大100件）
- CSVファイルはUTF-8エンコーディングです

## 📄 ライセンス

個人・研究用途で自由に使用できます♡

---

✨ 楽しい論文収集を！ ✨
