import pandas as pd

# Goが作ったCSVを読み込む
df = pd.read_csv("paper_box/face swap/2025/face swap_2025.csv")

# データを確認
print(df.head())
print(f"収集した論文数: {len(df)}")

# ここからNLP（自然言語処理）でタイトルの傾向分析などができます