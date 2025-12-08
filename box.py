import os
import csv
import shutil
import pandas as pd
from pathlib import Path

# paper_boxフォルダの作成
PAPER_BOX_DIR = "paper_box"
if not os.path.exists(PAPER_BOX_DIR):
    os.makedirs(PAPER_BOX_DIR)

# originalフォルダの作成
ORIGINAL_DIR = os.path.join(PAPER_BOX_DIR, "original")
if not os.path.exists(ORIGINAL_DIR):
    os.makedirs(ORIGINAL_DIR)

# 現在のディレクトリからCSVファイルを探す
csv_files = [f for f in os.listdir('.') if f.endswith('.csv')]

if not csv_files:
    print("❌ CSVファイルが見つかりませんでした...悲しい")
    exit(1)

print(f"📁 見つかったCSVファイル: {csv_files}ね♡")

# グローバルに見たURLを追跡（すべてのCSVファイル間で重複チェック）
global_seen_urls = set()

# 既存の paper_box フォルダから既に保存済みのURLを読み込む
def load_existing_urls():
    existing_urls = set()
    if os.path.exists(PAPER_BOX_DIR):
        for root, dirs, files in os.walk(PAPER_BOX_DIR):
            for file in files:
                if file.endswith('.csv'):
                    csv_path = os.path.join(root, file)
                    try:
                        with open(csv_path, 'r', encoding='utf-8') as f:
                            reader = csv.DictReader(f)
                            for row in reader:
                                url = row.get('URL', '')
                                if url:
                                    existing_urls.add(url)
                    except:
                        pass
    return existing_urls

# 既存のURLを読み込む
global_seen_urls = load_existing_urls()
print(f"📊 既存の論文: {len(global_seen_urls)}件を確認しました♡\n")

for csv_file in csv_files:
    print(f"\n🔄 処理中: {csv_file}...頑張ってるよ💪")
    
    try:
        with open(csv_file, 'r', encoding='utf-8') as f:
            reader = csv.DictReader(f)
            
            # Keywordごとにデータをグループ化
            keyword_groups = {}
            duplicate_count = 0  # 重複件数をカウント
            
            for row in reader:
                url = row.get('URL', '')
                # URLで重複チェック（グローバルで追跡）
                if url in global_seen_urls:
                    duplicate_count += 1
                    continue
                global_seen_urls.add(url)
                
                keyword = row.get('Keyword', 'Unknown')
                published = row.get('Published', '')
                # Published から年を抽出 (YYYY-MM-DDTHH:MM:SSZ の形式)
                year = published[:4] if len(published) >= 4 else 'Unknown'
                
                if keyword not in keyword_groups:
                    keyword_groups[keyword] = {}
                
                if year not in keyword_groups[keyword]:
                    keyword_groups[keyword][year] = []
                
                keyword_groups[keyword][year].append(row)
            
            # キーワード別にフォルダを作成してCSVを分割
            for keyword, year_groups in keyword_groups.items():
                # フォルダ名をクリーニング（ファイルシステムの無効文字を削除）
                folder_name = keyword.replace('/', '_').replace('\\', '_').replace(':', '_')
                folder_path = os.path.join(PAPER_BOX_DIR, folder_name)
                
                # キーワードフォルダを作成
                if not os.path.exists(folder_path):
                    os.makedirs(folder_path)
                
                # 年ごとにサブフォルダを作成
                for year, papers in year_groups.items():
                    year_folder = os.path.join(folder_path, year)
                    
                    if not os.path.exists(year_folder):
                        os.makedirs(year_folder)
                    
                    # 年別CSVファイルを作成or追記
                    output_csv = os.path.join(year_folder, f"{folder_name}_{year}.csv")
                    
                    # 既存ファイルがあれば読み込んで追記、なければ新規作成
                    if os.path.exists(output_csv):
                        existing_df = pd.read_csv(output_csv)
                        existing_urls = set(existing_df['URL'].dropna())
                        # 新規論文のみフィルタリング
                        new_papers = [p for p in papers if p.get('URL', '') not in existing_urls]
                        if new_papers:
                            with open(output_csv, 'a', newline='', encoding='utf-8') as out_f:
                                writer = csv.DictWriter(out_f, fieldnames=reader.fieldnames)
                                writer.writerows(new_papers)
                            print(f"   ✅ {folder_name}/{year}: {len(new_papers)}件の新規論文を追記しました💝 → {output_csv}")
                        else:
                            print(f"   ℹ️  {folder_name}/{year}: 新規論文なし（既に全て保存済み）")
                    else:
                        with open(output_csv, 'w', newline='', encoding='utf-8') as out_f:
                            writer = csv.DictWriter(out_f, fieldnames=reader.fieldnames)
                            writer.writeheader()
                            writer.writerows(papers)
                        print(f"   ✅ {folder_name}/{year}: {len(papers)}件の論文を整理しました💝 → {output_csv}")
        
        # 元のCSVファイルを paper_box/original へ移動
        dest_csv = os.path.join(ORIGINAL_DIR, csv_file)
        shutil.move(csv_file, dest_csv)
        print(f"   📂 {csv_file} をお片付けしました〜 → {dest_csv}")
        if duplicate_count > 0:
            print(f"   ⚠️  {duplicate_count}件の重複論文はスキップしちゃった♡")
        
    except Exception as e:
        print(f"   ❌ あ、エラーが出ちゃった...ごめんね: {e}")

print(f"\n✨ 完了しちゃった〜！論文は {PAPER_BOX_DIR} フォルダにきれいに整理されました♡")
