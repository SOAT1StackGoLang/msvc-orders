-- Drop the foreign key constraint from lanchonete_orders table
ALTER TABLE public.lanchonete_orders DROP CONSTRAINT fk_order_user_id;

-- Drop the lanchonete_users table
DROP TABLE IF EXISTS public.lanchonete_users;
