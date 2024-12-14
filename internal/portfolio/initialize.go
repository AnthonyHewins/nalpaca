package portfolio

import "context"

func (c *Controller) Init(ctx context.Context) error {
	_, err := c.nalpaca.GetAccount()
	if err != nil {
		c.logger.Error("failed fetching acct", "err", err)
		return err
	}

	return nil
}
