INSERT INTO budget_item(id, code, name, level, accumulate, parent_id) VALUES
('439082ad-f1bd-4228-91f2-8e744894ffdc', '500', 'Costo Directo', 1, true, null),
('b4b2e4e4-f22d-402e-9ab5-1d59347cbfcb', '500.1', 'Obra Gruesa', 2, false, '439082ad-f1bd-4228-91f2-8e744894ffdc'),
('420f8bb3-bc8e-4564-be99-75cd7c1a6ff8', '200', 'Costo Urbanizacion', 1, true, null),
('9abc2426-a92b-46ef-b074-ddbc8ee2df1a', '200.1', 'Adoquin', 2, false, '420f8bb3-bc8e-4564-be99-75cd7c1a6ff8');
