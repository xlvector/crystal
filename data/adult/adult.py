head = ['label', '*age', '#workclass', '*fnlwgt', '#education', '*education-num', '#marital-status', '#occupation', '#relationship', '#race', '#sex', '*capital-gain', '*capital-loss', '*hours-per-week', '#native-country']

print '\t'.join(head)
for line in file('adult.data'):
    tks = line.strip().split(',')
    tks = [x.strip() for x in tks]
    category = tks[-1]
    print category + '\t' + '\t'.join(tks[:-1])

