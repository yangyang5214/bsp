import os
import fnmatch
import docx
from typing import List
import re
import pandas as pd
import numpy as np

np.float = np.float64


class Customer:
    def __init__(self, content: List[str]):
        self.content = content
        self.name = ''
        self.age = ''
        self.last_date = ''
        self.desc = ''  # 症状
        self.guide = ''  # 药方

    def split(self, v: str) -> List[str]:
        arrs = v.split(' ')
        result = [item for item in arrs if item != ""]

        if len(result) == 1:
            result = result + ['', '']
        elif len(result) == 2:
            if '.' in result[1]:
                result.append(result[1])
                result[1] = ''
            else:
                result.append('')
        return result

    def parser(self) -> List[str]:
        if not self.content:
            return []

        first_line = self.content[0]
        arrs = self.split(first_line)

        self.name = arrs[0]
        self.age = arrs[1]
        self.last_date = arrs[2]

        index_flag = 1
        for i in range(1, len(self.content)):
            item = self.content[i]
            if re.match(r'.*\d+.*', item):
                index_flag = i
                break

        if index_flag > 1:
            self.desc = ''.join(self.content[1:index_flag])
            self.guide = ''.join(self.content[index_flag + 1:])
        else:
            self.guide = ''.join(self.content[index_flag:])

        return [
            self.name,
            self.age,
            self.last_date,
            self.desc,
            self.guide,
            "\n".join(self.content),
        ]


def find_all_file(dir_name: str, extension: str):
    for root, _, files in os.walk(dir_name):
        for basename in files:
            if fnmatch.fnmatch(basename, f'*{extension}'):
                filename = os.path.join(root, basename)
                yield filename


def process_single(file_name: str) -> list:
    print(f'start read file:{file_name}')
    doc = docx.Document(file_name)
    full_text = []

    sub_text = []
    for para in doc.paragraphs:
        text = para.text
        sub_text.append(text)

        if text == '' and len(sub_text) != 0 and sub_text[0] != '':
            full_text.append(sub_text)
            sub_text = []
    return full_text


def split_keyword(content: str) -> List:
    r = []
    arrs = content.replace('，', ' ').split(" ")
    for item in arrs:
        if item == '':
            continue

        r.append(item)
    return r


def main():
    keys1 = []
    keys2 = []
    for name in find_all_file("/Users/beer/Downloads/hdd", "docx"):
        rs = process_single(name)
        for r in rs:
            c = Customer(r)
            c.parser()
            # result.append(['/'.join(name.split("/")[-2:])] + c.parser())
            keys1 = keys1 + split_keyword(c.desc)
            keys2 = keys2 + split_keyword(c.guide)

    # df = pd.DataFrame(result, columns=['文件名', '客户', '年龄', '最后看病日期', '症状', '药单', '原始文本'])
    df = pd.DataFrame([
        keys1,
    ]).melt()
    with pd.ExcelWriter('/tmp/keyword_1.xlsx', engine='openpyxl') as writer:
        df.to_excel(writer, index=False)

    df = pd.DataFrame([
        keys2,
    ]).melt()
    with pd.ExcelWriter('/tmp/keyword_2.xlsx', engine='openpyxl') as writer:
        df.to_excel(writer, index=False)




if __name__ == '__main__':
    main()
