;; if x is false, then y is not even evaluated
(defun and (x y) 
  (cond
   (x (cond
       (y 't)
       ('t ())))
   ('t '())))
