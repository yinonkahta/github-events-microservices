let error = true

let res = [
    db.events.drop(),
    db.events.createIndex({ 'event.created_at': -1 }),
    db.repos.createIndex({ 'last_updated_at': -1 }),
    db.users.createIndex({ 'last_updated_at': -1 }),
]

printjson(res)

if (error) {
    print('Error, exiting')
    quit(1)
}