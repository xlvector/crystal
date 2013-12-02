crystal
=======

Web based Data Visualization Tools

## How to Run

If you have a dataset called "adult",

  cd data
  mkdir adult
  mv path_to_you_dataset adult/data.tsv
  
Then

  go build crystal.go
  nohup ./crystal > log &

Then, you can view it through http://localhost:8080/html/dim1.html

## Dataset Format

Following is a sample dataset

  label   age     [workclass]     fnlwgt  [education]
  <=50K   39      State-gov       77516   Bachelors
  <=50K   50      Self-emp-not-inc        83311   Bachelors
  >50K    40      Private 121772  Assoc-voc
  
Here, label is the label of sample. age, [workclass], fnlwgt, [education] are features. Age, fnlwgt are continuous features and [workclass], [education] are discrete features.
