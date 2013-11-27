import sys, json

fname, = sys.argv[1:]

data = []
head = []

n = 0
for line in file(fname):
    if len(line.strip()) == 0:
        continue
    n += 1
    if n == 1:
        head = line.strip().split('\t')[1:]
    else:
        tks = line.strip().split('\t')
        data.append((tks[0], tks[1:]))

stat = {}
col_values = {}
for label, sample in data:
    for i in range(len(head)):
        if head[i][0] != '#':
            continue
        col = head[i][1:]
        if col not in stat:
            stat[col] = {}
            col_values[col] = set()
        if label not in stat[col]:
            stat[col][label] = {}
        if sample[i] not in stat[col][label]:
            stat[col][label][sample[i]] = 0
            col_values[col].add(sample[i])
        stat[col][label][sample[i]] += 1

prefix = fname.split('.')[0]
files = []
for col, dis in stat.items():
    key_sum = {}
    records = []
    for label, values in dis.items():
        record = {}
        record["key"] = col + ': ' + label
        record["values"] = []
        for key in col_values[col]:
            count = 0
            if key in values:
                count = values[key]
            record["values"].append({"x": key, "y": count})
            if key not in key_sum:
                key_sum[key] = 0.0
            key_sum[key] += float(count)
        records.append(record)
    with open(prefix + '_' + col + '_stacked.json', 'w') as sw:
        json.dump(records, sw, indent=4)

    records = []
    for label, values in dis.items():
        record = {}
        record["key"] = col + ': ' + label
        record["values"] = []
        for key in col_values[col]:
            count = 0.0
            if key in values:
                count = values[key]
            record["values"].append({"x": key, "y": count / key_sum[key]})
        records.append(record)
    with open(prefix + '_' + col + '_100stacked.json', 'w') as sw:
        json.dump(records, sw, indent=4)

    files.append({"name": col, "charts" : [prefix + '_' + col + '_stacked.json', prefix + '_' + col + '_100stacked.json']})

with open('list.json', 'w') as sw:
    json.dump(files, sw, indent=4)
