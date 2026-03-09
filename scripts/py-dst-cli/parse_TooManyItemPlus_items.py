import os
import re
import json

# 定义修饰词规则
SPICE_TRANSLATIONS = {
    "_SPICE_GARLIC": "蒜",
    "_SPICE_CHILI": "辣",
    "_SPICE_SUGAR": "甜",
    "_SPICE_SALT": "盐",
}

def parse_po_file(po_file):
    """解析 .po 文件，从第 9 行开始，返回仅包含 STRINGS.NAMES 的 msgctxt 与 msgstr 的映射"""
    translations = {}
    with open(po_file, "r", encoding="utf-8") as f:
        lines = f.readlines()

    # 从第 9 行开始解析
    lines = lines[8:]

    msgctxt, msgid, msgstr = None, None, None
    for line in lines:
        line = line.strip()
        if line.startswith("msgctxt") and "STRINGS.NAMES." in line:
            msgctxt = re.search(r'"STRINGS\.NAMES\.(.+)"', line).group(1)  # 提取物品名
        elif line.startswith("msgid"):
            msgid = re.search(r'"(.*)"', line).group(1)  # 捕获空字符串
        elif line.startswith("msgstr"):
            msgstr = re.search(r'"(.*)"', line).group(1)  # 捕获空字符串

            # 只有在 msgctxt、msgid 和 msgstr 都非空时才保存
            if msgctxt and msgid and msgstr:
                translations[msgctxt] = msgstr

            # 重置状态
            msgctxt, msgid, msgstr = None, None, None
    return translations

def apply_spice_rule(item, base_translation):
    print(item)
    """根据规则生成修饰后的翻译"""
    for suffix, spice_translation in SPICE_TRANSLATIONS.items():
        if item.endswith(suffix):
            base_item = item.replace(suffix, "")
            if base_item in base_translation:
                return f"{base_translation[base_item]}-{spice_translation}"
    return "未翻译"

def generate_translations(input_folder, po_file, output_file):
    """生成带翻译的 JSON 数据"""
    # 解析 .po 文件
    translations = parse_po_file(po_file)

    # 初始化结果字典
    result = {}

    # 遍历输入文件夹，处理分类数据
    for filename in os.listdir(input_folder):
        if filename.endswith(".lua"):
            category_name = os.path.splitext(filename)[0]
            file_path = os.path.join(input_folder, filename)
            with open(file_path, "r", encoding="utf-8") as f:
                content = f.read()
            if "return" in content:
                content = content.split("return", 1)[1].strip().strip("{}")
                items = [item.strip().strip('"') for item in content.split(",")]
                result[category_name] = {
                    item: translations.get(item.upper(), apply_spice_rule(item.upper(), translations))
                    for item in items
                }

    # 保存结果为 JSON 文件
    with open(output_file, "w", encoding="utf-8") as f:
        json.dump(result, f, ensure_ascii=False, indent=4)
    print(f"翻译结果已保存到 {output_file}")


# 设置文件路径
input_folder = "./scripts/TMIP/list"  # 替换为文件夹路径
po_file = "./chinese_s.po"           # 替换为 .po 文件路径
output_file = "tooManyItemPlus.json"     # 输出文件路径

# 运行生成函数
generate_translations(input_folder, po_file, output_file)
