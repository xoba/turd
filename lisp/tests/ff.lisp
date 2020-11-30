;; create a function in a function
;; TODO: need to evaluate x, but then quote it!
((lambda (x)
   (list 'lambda '(y) (list 'cons (list 'quote x) 'y)))
 '5)
