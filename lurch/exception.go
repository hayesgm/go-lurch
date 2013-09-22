package lurch

import (
  "labix.org/v2/mgo"
  "labix.org/v2/mgo/bson"
  "time"
  "log"
  "hash/fnv"
  "path"
)

type Exception struct {
  Id bson.ObjectId "_id,omitempty"
  UniqueId uint32
  Name string
  Message string
  Status string
  Comments []string
  Instances []time.Time
  Traces []*Trace `bson:"Trace",omitempty`
}

// Create a new Exception
// This will generate a UniqueId for the error
func generateUniqueId(name, message string, backtrace []string) (uniqueId uint32) {

  id := path.Join(name, message)
  
  // Add first line of backtrace if available
  if len(backtrace) > 0 {
    id = path.Join(id, backtrace[0])
  }

  hash := fnv.New32a()
  hash.Write([]byte(id))

  return hash.Sum32()
}

func NewException(name, message string, backtrace []string) (exception *Exception, err error) {
  exception = &Exception{UniqueId: generateUniqueId(name, message, backtrace), Name: name, Message: message}

  return
}

// Find an Exception by UniqueId from Mongo
func FindException(uniqueId uint32, db *mgo.Database) (exception *Exception, err error) {
  
  exception = &Exception{}
  err = db.C("exceptions").Find(bson.M{"uniqueid": uniqueId}).One(&exception)
  if err != nil {
    return nil, err
  }

  return
}

// Find an Exception by UniqueId from Mongo
func GetExceptions(db *mgo.Database) (exceptions []*Exception, err error) {
  
  iter := db.C("exceptions").Find(nil).Iter()
  err = iter.All(&exceptions)
  if err != nil {
    return
  }

  return
}

func (exception *Exception) Insert(db *mgo.Database) (err error) {
  err = db.C("exceptions").Insert(exception)
  if err != nil {
    return err // That's a failure
  }

  return
}

// Upsert will insert a new Exception if nothing matching uniqueid exists
// Then, it will add to instances/traces of the Exception
func (exception *Exception) Upsert(trace *Trace, db *mgo.Database) (err error) {

  log.Printf("Upserting: %#v\n", exception)

  if err = db.C("exceptions").EnsureIndex(mgo.Index{Key: []string{"uniqueid"}, Unique: true}); err != nil {
    return
  }
  
  // First, we'll try to insert
  // Due to the unique constraint on UniqueID, this will fail if the record exists already
  err = db.C("exceptions").Insert(&exception)
  if err != nil && !mgo.IsDup(err) { // Swallow duplicate error
    return
  }

  // Now, we're going to push the instance to mark when this happened
  err = db.C("exceptions").Update(bson.M{"uniqueid": exception.UniqueId}, bson.M{
    "$push": bson.M{"instances": time.Now() },
  })

  if err != nil {
    return err
  }

  // Finally, we'll add the trace (for now, unlimited)
  err = db.C("exceptions").Update(bson.M{"uniqueid": exception.UniqueId}, bson.M{
    "$push": bson.M{"traces": &trace },
  })

  if err != nil {
    return err
  }

  return
}