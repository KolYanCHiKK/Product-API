package repository

const createProduct = `
	INSERT INTO products(
		name, 
	   	price, 
	    quantity,
	    descriptions,
	    images                 
	)
	VALUES($1, $2, $3, $4, $5)
	RETURNING 
		product_id,
		name, 
	   	price, 
	    quantity,
	    descriptions,
	    images,  
		created_at,
		updated_at
`

const getAllProducts = `
	SELECT
	    product_id,
		name, 
	   	price, 
	    quantity,
	    descriptions,
	    images,  
		created_at,
		updated_at
	FROM products
	ORDER BY %s product_id
	LIMIT $1
	OFFSET $2
`

const countAllProductsRows = `
	SELECT COUNT(*)
	FROM products
`

const getProductByID = `
	SELECT *
	FROM products p
	where p.product_id = $1
`

const updateProduct = `
	UPDATE products
	SET 
	    name = $1,
		price = $2,
		quantity = $3,
		descriptions = $4,
	    images[1] =  $5
	WHERE product_id = $6
	RETURNING *
`

const patchProduct = `
	UPDATE products
	SET %s
	WHERE product_id = $%d
	RETURNING *
`

const deleteProduct = `
	DELETE FROM products
	where %s
	RETURNING *
`
