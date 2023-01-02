package product

type Service interface {
	Create(product InputProduct) (Product, error)
	GetAll() ([]Product, error)
	GetById(id int) (Product, error)
	Update(id int, inputTask InputProduct) (Product, error)
	Delete(id int) error
}

type service struct {
	repository Repository
}

func NewService(r Repository) *service {
	return &service{r}
}

func (s *service) Create(input InputProduct) (Product, error) {
	var task Product
	task.Name = input.Name
	task.Price = input.Price

	newTask, err := s.repository.Create(task)
	if err != nil {
		return task, err
	}

	return newTask, nil
}

func (s *service) GetAll() ([]Product, error) {
	tasks, err := s.repository.GetAll()
	if err != nil {
		return tasks, err
	}

	return tasks, nil
}

func (s *service) GetById(id int) (Product, error) {
	task, err := s.repository.GetById(id)
	if err != nil {
		return task, err
	}

	return task, nil
}

func (s *service) Update(id int, input InputProduct) (Product, error) {
	uTask, err := s.repository.Update(id, input)
	if err != nil {
		return uTask, err
	}

	return uTask, nil
}

func (s *service) Delete(id int) error {
	err := s.repository.Delete(id)
	if err != nil {
		return err
	}

	return nil
}
