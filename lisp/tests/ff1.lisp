;; try to use a function created within a function
(funcall ((lambda (x)
   (list 'lambda '(y) (list 'cons (list 'quote x) 'y)))
 '5)
'(1 2 3))
