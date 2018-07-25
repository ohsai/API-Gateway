import Data.List
import System.Random
import Control.Monad (replicateM)

main = do
    let a = [[1,2],[3,4]]
    let b = [[5,6,7],[7,8,9]]
    let c = [[1.0,2.0],[3.0,4.0]]
    let d = [[-1.0,-2.0],[1.0,2.0]]
    --printm $ mmult a b
    --printm $ unsqueeze (3,2) [1.0,2.0,3.0]
    --print $ squeeze [[1],[2],[3],[4],[5]]
    {-
    let (out, cache) = affine_layer a [1.0,1.0] c
    let (dx, dw, db) = affine_layer_backward cache out
    printm dx
    printm dw
    print db
    -}
    --printm $ mmult_elem a (transpose c)
    printm $ snd (softmax_loss d [0,1])
    printm $ indexing_one 10 [0,1,2,3,4,5,6,7,8,9]
    seed <- newStdGen
    let layer_print (x,y) = do 
        printm x
        print y
        putStrLn ""
    mapM_ layer_print $ multi_layer_net_initialize [5,5,5] seed 
    train_default [] [] [] [] []
    brl <- boundedrandomList (0,10) 10
    print brl

mmult :: (Num a) => [[a]] -> [[a]] -> [[a]]
mmult a b = if snd (dimension a) == fst (dimension b)
    then [ [ sum $ zipWith (*) ar bc | bc <- (transpose b) ] | ar <- a ]
    else error "matrix format does not comply for mmult"
madd :: (Num a) => [[a]] -> [[a]] -> [[a]]
madd a b = if dimension a == dimension b
    then zipWith(zipWith (+) ) a b
    else error "matrix format does not comply for madd"
mzeros :: (Num a ) => Int -> Int -> [[a]]
mzeros x y = replicate x $ replicate y 0
maddinv :: (Num a) => [[a]] -> [[a]]
maddinv a = map (map negate) a
smmult :: (Num a) => a -> [[a]] -> [[a]]
smmult a b = map (map (*a) ) b
mones :: (Num a) => Int -> Int -> [[a]]
mones x y = replicate x $ replicate y 1
reduce_sum :: (Num a) => [[a]] -> Int -> [a]
reduce_sum a b 
    | null (head a) = []
    | b == 1 = map (sum) a
    | b == 0 = (sum (map head a)) : (reduce_sum (map tail a) b)
    | otherwise = error "reduce_sum axis must be 0 or 1"
dimension :: (Num a) => [[a]] -> (Int,Int)
dimension p = (length p, length (head p))
printm :: (Num a, Show a) => [[a]] -> IO ()
printm p = do
    let (x,y) = dimension p
    putStrLn $ (show x) ++ " x " ++ (show y) ++ " matrix"
    mapM_ print (p)
    putStrLn ""
unsqueeze :: (Num a) => (Int, Int) -> [a] -> [[a]]
unsqueeze (a,b) x 
    | length x == b = replicate a x
    | length x == a = transpose $ replicate b x 
    | otherwise = error $ "Unable to unsqueeze length "++(show (length x))  ++ " to "++ (show (a,b))
squeeze :: (Num a ) => [[a]] -> [a]
squeeze x 
    | fst (dimension x) == 1 =  head x
    | snd (dimension x) == 1 = head (transpose x)
    | otherwise = error "Unable to squeeze 2d array with dimensions like that"
mmult_elem :: (Num a) => [[a]] -> [[a]] -> [[a]]
mmult_elem a b
    | dimension a == dimension b = zipWith (zipWith (*) ) a b
    | otherwise = error "Unable to elementwise multiply two matrices"
mdiv_elem :: (Fractional a) => [[a]] -> [[a]] -> [[a]]
mdiv_elem a b
    | dimension a == dimension b = zipWith (zipWith (/) ) a b
    | otherwise = error "Unable to elementwise multiply two matrices"

-- transpose is on Data.List already

affine_layer :: (Num a) => [[a]] -> [a] -> [[a]] -> ([[a]],([[a]],[[a]],[a])) 
affine_layer weight bias input = (out, (input, weight, bias))
    where
        out = madd (mmult input weight ) (unsqueeze res_dim bias)
        res_dim = (fst (dimension input), snd (dimension weight))

affine_layer_backward :: (Num a) => ([[a]],[[a]],[a])-> [[a]] -> ([[a]],[[a]],[a])  
affine_layer_backward (x,w,b) dout = (dx, dw, db)
    where
        dx = mmult dout (transpose w)
        dw = mmult (transpose x) dout
        db = squeeze(mmult (mones 1 (snd (dimension x))) dout)

relu_forward :: (Num a, Ord a) => [[a]] -> ([[a]], [[a]])
relu_forward x = (map (map (max 0)) x , x)

relu_backward :: (Num a, Ord a) => [[a]] -> [[a]] -> [[a]]
relu_backward x dout = mmult_elem (map (map (\y -> if y <= 0 then 0 else 1)) x) dout

softmax_loss :: ( Ord a,Floating a) => [[a]] -> [Int] -> (a,[[a]])
softmax_loss x y = (loss,dx)
    where
        probs = map (map (exp) ) $madd x $maddinv $unsqueeze (dimension x) (map maximum x)
        probs_2 = mdiv_elem probs (unsqueeze (dimension probs) (reduce_sum probs 1))
        nN = fst (dimension x)
        loss = (-1.0 /(fromIntegral nN)) * sum (map log (zipWith (!!) probs_2 y))
        dx = map (map (/(fromIntegral nN))) $madd probs_2 $maddinv $indexing_one (snd (dimension probs_2)) y 
        
indexing_one :: (Num a) => Int -> [Int] -> [[a]]
indexing_one max_index label = map (\x -> (zerolist x) ++ [1] ++ (zerolist (max_index - 1 - x))) label 
    where
        zerolist n = replicate n 0 

multi_layer_net :: (Floating a, Ord a) => a ->[([[a]],[a])]->[[a]] ->[Int] ->  
    (a,[([[a]],[a])])
multi_layer_net reg model input y = (loss, grad)
    where
        (scores,cache) = mapAccumL (\(prevout) -> (\(weight,bias) -> 
            let (curout, newstate) = affine_layer weight bias prevout 
            in
                let (curout_relu, newstate_relu) = relu_forward curout
                in
                (curout_relu, (newstate_relu , newstate) )
                )) input model 
        {-
        (scores,cache) = foldl (\(prevout,state_cache) -> (\(weight,bias) -> 
            let (curout, newstate) = affine_layer weight bias prevout 
            in
                let (curout_relu, newstate_relu) = relu_forward curout
                in
                (curout_relu, newstate_relu : newstate : state_cache )
                )) (X,[]) model 
        -}
        (data_loss , dx )= softmax_loss scores y 
        loss = 0.5 * reg * (foldl (\reg_loss -> (\(weight, bias) -> 
            (reg_loss + sum (reduce_sum (mmult_elem weight weight) 1))
            )) 0 model) + data_loss
        
        (_,grad_rev,_) = foldr (\(weight, bias) -> (\(dprev_out , grad_cache, cache_left) -> 
            let (dout, dw,db) = affine_layer_backward (snd (last cache_left)) dprev_out
            in
                let (dout_relu) = relu_backward (fst (last cache_left)) dout
                in
                (dout, (dw, db): grad_cache,(init cache_left))
             )) (dx, [],cache) model 
        grad = reverse grad_rev

multi_layer_net_initialize :: ( Num b) => [Int] -> StdGen -> [([[Float]],[b])]
multi_layer_net_initialize layer_dimensions seed = model
    where
        model = pl2lp wl bl 
        wl = map  weight_initialization_method (pl2lp (init layer_dimensions) (tail layer_dimensions))
        bl = map bias_initialization_method (tail layer_dimensions)
        pl2lp l1 l2 = if length l1 == length l2 
            then zipWith (\x -> \y -> (x,y)) l1 l2
            else error "trying 2 pairize unequal length lists"
        weight_initialization_method = \x -> smmult (0.00001) ((uncurry randomMatrix x) seed)
        bias_initialization_method = \x -> replicate x 0 
randomMatrix :: Int -> Int -> StdGen -> [[Float]] --range [0,1)
randomMatrix n m seed = (splitEvery m)$ randomlist (n*m) seed
randomlist ::  Int -> StdGen -> [Float]
randomlist n = take n . unfoldr (Just . random)

splitEvery :: Int -> [a] -> [[a]]
splitEvery _ [] = []
splitEvery n xs = as : splitEvery n bs 
  where (as,bs) = splitAt n xs

stochastic_gradient :: (Num a) =>a -> [([[a]],[a])] -> [([[a]],[a])]-> [([[a]],[a])]
stochastic_gradient learning_rate model grad = updated_model
    where 
        updated_model = zipWith (update) model grad
        update (weight, bias) (weight_grad, bias_grad) = ((madd weight (smmult (-1 * learning_rate) weight_grad)),(zipWith (-) bias (map ((*) learning_rate) bias_grad) )) 

boundedrandomList :: (Int, Int) -> Int -> IO([Int])
boundedrandomList (lower_bound, greater_bound) n = replicateM n $ randomRIO (lower_bound,greater_bound)

train :: (Num a) => Float -> Float -> Float -> Float -> Int -> Int -> Bool -> 
    String -> -- update
    ([[a]] -> [Int] -> (a,[[a]]) )->  --loss function
    [([[a]],[a])] -> [[a]] -> [a] -> [[a]] -> [a] -> IO ([([[a]],[a])])
train reg learning_rate momentum learning_rate_decay num_epochs batch_size verbose updater loss_function model x y x_val y_val = do
    let number_of_data = fst $dimension x
    let number_of_val_data = fst $dimension x_val
    let iterations_per_epoch = number_of_data / batch_size
    train_step steps_left model_input = 
        if steps_left == 0 then model_input else  train_step (steps_left -1) updated_model_input
        where
            chosen_indices <- boundedrandomList (0,number_of_data-1) batch_size
            x_batch = map (x !!) chosen_indices
            y_batch = map (y !!) chosen_indices
            (train_loss, train_grad) = multi_layer_net reg model x_batch y_batch
            updated_model_input = stochastic_gradient learning_rate model train_grad
    train_step (iterations_per_epoch * num_epochs) model

    --let (validation_loss , _) = multi_layer_net reg model x_val_batch y_val_batch
    print("Hello World!")

train_default = train 0.0 1e-2 0 0.95 30 100 True "sgd" softmax_loss 


    
    

