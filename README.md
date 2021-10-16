# ToDoList-API
todo list api where you can retrieve, update, insert values 


## SignUp
pass username and password as form value to path below
*/api/signup* 
key is : 
  username -- enter preferred username
  password -- enter preferred username

## Login
pass exsisting username and password as form value to path below
*/api/login*                       
key is : 
  username -- enter preferred username
  password -- enter preferred username

## get all todolists
after you have logged in or signed up, you can get your todolists with path below
*/api/getTodoList*

## get completed items
*/api/getCompletedItems*

## get incomplete items
*/api/getIncompleteItems*

## insert new item
pass title, desription, completed key & value as form value with path below (POST)
*/api/insertItem*

## update items
currently you can only update the "completed" property of particular todo.
(will be adding feature to update other properties later)
pass either "true" or "false" as form value (key is "completed") to path below (POST)
*/api/updateItem?id=<id>*


