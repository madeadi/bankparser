# Bank Parser

To parse the data from the bank to be imported to GNU Cash. 


## Providers

Right now, this program can read the data from DBS and Revolut. 

- `revolut`: download the CSV from the revolut app. 
- `dbs`: copy the PDF table and paste it into a text file. 

## How to use

```
go run main.go -in ./samples/revolut.txt -out ./out/dbs.csv -provider {choose one: revolut|dbs}
```

## Importing transactions to GNU Cash

Please refer to this [YouTube Video](https://www.youtube.com/watch?v=h6TIviAn6kg)