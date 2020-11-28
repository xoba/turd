;; this doesn't seem to work
(defun y (f)
  ((lambda (x) (f (x x)))
   (lambda (x) (f (x x)))))
