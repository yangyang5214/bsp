import pandas as pd
import numpy as np
import re

np.float = np.float64


def parser_excel():
    df = pd.read_excel("/tmp/1.xlsx", sheet_name='bk', engine='openpyxl')
    values = df.values

    results = []

    title = ''
    for item in values:
        v = str(item[1])
        if is_title(v):
            title = v
            continue
        results.append([title, v])

    return pd.DataFrame(results, columns=['title', 'text'])


def is_title(v: str) -> bool:
    if not re.match(r'^\d+\s.*', v):
        return False
    if len(v) > 4 and '.' in v:
        return False
    return True


def main():
    df = parser_excel()
    with pd.ExcelWriter('/tmp/result.xlsx', engine='openpyxl') as writer:
        df.to_excel(writer)


if __name__ == '__main__':
    main()
