db.getSiblingDB('admin').createUser(
        {
            user: "flamingo-group",
            pwd: "flamingo-group",
            roles: [
                {
                    role: "readWrite",
                    db: "flamingo-group"
                },
                {
                    role: "read",
                    db: "admin"
                }
            ]
        }
);

db.getSiblingDB('flamingo-group').createUser(
        {
            user: "flamingo-group",
            pwd: "flamingo-group",
            roles: [
                {
                    role: "readWrite",
                    db: "flamingo-group"
                },
                {
                    role: "read",
                    db: "admin"
                }
            ]
        }
);
