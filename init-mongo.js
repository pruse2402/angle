db.createUser(
    {
        user : "angle",
        pwd : "angle",
        roles : [
            {
                role : "readwrite",
                db   : "angle"
            }
        ]
    }
)