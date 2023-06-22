Next, I will outline some details that the code may include, which, due to time constraints, I was unable to implement in the most suitable manner.

Initially, concerning the database (DB), one possible approach would be to create interfaces for each schema. This way, we can better separate the functionalities associated with each table.

On another note, regarding testing, what I did was perform integration tests that cover the entire flow from the handlers to the service, testing each of the possible scenarios. Unit testing, on the other hand, focuses on everything related to the database repository.

I made this decision because if we had a complete integration test including the handler and the DB, any future changes to the DB would require fixing all the tests from the handlers. With this approach, we isolate the database unit testing and keep it separate. In case of a change, it remains transparent to the service without anyone noticing.

The project currently has approximately 90% code coverage.

Lastly, I decided to allow multiple investors to contribute funds towards purchasing invoices. Only when the requested amount for an invoice is reached do we proceed with the actual purchase. This way, an invoice can have multiple investors.

Certainly, there are many validations that are missing and would make the code more robust. Specifically, error handling in case of failure during an operation and rollback of executed operations up to the point of failure are not considered in this application.

This MVP (Minimum Viable Product) adheres strictly to the requirements specified, but I believe it is a highly scalable application that can evolve according to our preferences.