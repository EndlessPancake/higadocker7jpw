from pymongo import MongoClient

class MongoSample(object):

    def __init__(self, name):
        self.client = MongoClient()
        self.db = self.client[name]
    def collection_names(self):
        return self.db.collection_names()
    def list_collection_names(self):
        return self.db.list_collection_names()
    def main(self):
        print(mongo.list_collection_names())

mongo = MongoSample('test')
if __name__ == "__main__":
   main()
