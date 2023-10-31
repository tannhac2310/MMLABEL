db.getSiblingDB('admin').createUser(
        {
            user: "mm-printing",
            pwd: "mm-printing",
            roles: [
                {
                    role: "readWrite",
                    db: "mm-printing"
                },
                {
                    role: "read",
                    db: "admin"
                }
            ]
        }
);

db.getSiblingDB('mm-printing').createUser(
        {
            user: "mm-printing",
            pwd: "mm-printing",
            roles: [
                {
                    role: "readWrite",
                    db: "mm-printing"
                },
                {
                    role: "read",
                    db: "admin"
                }
            ]
        }
);
