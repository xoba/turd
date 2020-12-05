;; if x is false, then y is not even evaluated.
;; on the other hand, in calling "and" via lambda,
;; both x and y are actually evaluated.
(defun and (x y) 
  (cond
   (x (cond
       (y 't)
       ('t ())))
   ('t '())))
